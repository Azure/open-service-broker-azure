package postgresql

import (
	"context"
	"fmt"

	postgresSDK "github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-12-01/postgresql" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

const (
	enabledARMString  = "Enabled"
	disabledARMString = "Disabled"
)

func getAvailableServerName(
	ctx context.Context,
	checkNameAvailabilityClient postgresSDK.CheckNameAvailabilityClient,
) (string, error) {
	for {
		serverName := uuid.NewV4().String()
		available, err := isServerNameAvailable(
			ctx,
			serverName,
			checkNameAvailabilityClient,
		)
		if err != nil {
			return "", fmt.Errorf(
				"error determining server name availability: %s",
				err,
			)
		}
		if available {
			return serverName, nil
		}
	}
}

func isServerNameAvailable(
	ctx context.Context,
	serverName string,
	checkNameAvailabilityClient postgresSDK.CheckNameAvailabilityClient,
) (bool, error) {
	nameAvailability, err := checkNameAvailabilityClient.Execute(
		ctx,
		postgresSDK.NameAvailabilityRequest{
			Name: &serverName,
		},
	)
	if err != nil {
		return false, err
	}

	if *nameAvailability.NameAvailable {
		return true, nil
	}
	return false, nil
}

// generateDBMSInstanceDetail will read information
// from instance provision parameters, and generate
// a dbmsInstanceDetail. This method is expected to
// be invoked by preProvision step of all-in-one and
// dbms.
func generateDBMSInstanceDetails(
	ctx context.Context,
	instance service.Instance,
	checkNameAvailabilityClient postgresSDK.CheckNameAvailabilityClient,
) (*dbmsInstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Determine server name. If specified,
	// check availability; else, generate one.
	pp := instance.ProvisioningParameters
	serverName := pp.GetString("serverName")
	if serverName != "" {
		available, err := isServerNameAvailable(
			ctx,
			serverName,
			checkNameAvailabilityClient,
		)
		if err != nil {
			return nil, err
		} else if !available {
			return nil, fmt.Errorf("server name %s is already in use", serverName)
		}
	} else {
		var err error
		serverName, err = getAvailableServerName(
			ctx,
			checkNameAvailabilityClient,
		)
		if err != nil {
			return nil, err
		}
	}

	// Determine administratorLogin. If specified,
	// use it; else, use default value "postgres".
	adminAccountSettings := pp.GetObject("adminAccountSettings")
	adminUsername := adminAccountSettings.GetString("adminUsername")
	if adminUsername == "" {
		adminUsername = "postgres"
	}
	// Determine AdministratorLoginPassword. If specified,
	// use it; else, generate one.
	adminPassword := adminAccountSettings.GetString("adminPassword")
	if adminPassword == "" {
		adminPassword = generate.NewPassword()
	}
	return &dbmsInstanceDetails{
		ARMDeploymentName:          uuid.NewV4().String(),
		ServerName:                 serverName,
		AdministratorLogin:         adminUsername,
		AdministratorLoginPassword: service.SecureString(adminPassword),
	}, nil
}

func setupDatabase(
	enforceSSL bool,
	administratorLogin string,
	serverName string,
	administratorLoginPassword string,
	fullyQualifiedDomainName string,
	dbName string,
) error {
	db, err := getDBConnection(
		enforceSSL,
		administratorLogin,
		serverName,
		administratorLoginPassword,
		fullyQualifiedDomainName,
		primaryDB,
	)
	if err != nil {
		return err
	}
	defer db.Close() // nolint: errcheck

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %s", err)
	}
	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				log.WithField("error", err).Error("error rolling back transaction")
			}
		}
	}()
	if _, err = tx.Exec(
		fmt.Sprintf(`create role "%s"`, dbName),
	); err != nil {
		return fmt.Errorf(`error creating role "%s": %s`, dbName, err)
	}
	// Azure will automatically create a role having name of
	// postgreSQL server admin user.
	// Here the purpose of granting one role to another role is
	// to make admin role to be a member of created db role.
	// Please see: https://www.postgresql.org/docs/10/role-membership.html
	if _, err = tx.Exec(
		fmt.Sprintf(`grant "%s" to "%s"`, dbName, administratorLogin),
	); err != nil {
		return fmt.Errorf(
			`error adding role "%s" to role "%s": %s`,
			dbName,
			administratorLogin,
			err,
		)
	}
	if _, err = tx.Exec(
		fmt.Sprintf(
			`alter database "%s" owner to "%s"`,
			dbName,
			dbName,
		),
	); err != nil {
		return fmt.Errorf(
			`error updating database owner"%s": %s`,
			dbName,
			err,
		)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %s", err)
	}

	return nil
}

func createExtensions(
	enforceSSL bool,
	administratorLogin string,
	serverName string,
	administratorLoginPassword string,
	fullyQualifiedDomainName string,
	dbName string,
	extensions []string,
) error {
	db, err := getDBConnection(
		enforceSSL,
		administratorLogin,
		serverName,
		administratorLoginPassword,
		fullyQualifiedDomainName,
		dbName,
	)
	if err != nil {
		return err
	}
	defer db.Close() // nolint: errcheck

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %s", err)
	}
	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				log.WithField("error", err).Error("error rolling back transaction")
			}
		}
	}()
	for _, extension := range extensions {
		if _, err = tx.Exec(
			fmt.Sprintf(`create extension "%s"`, extension),
		); err != nil {
			return fmt.Errorf(
				`error creating extension "%s": %s`,
				extension,
				err,
			)
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %s", err)
	}
	return nil
}

func buildGoTemplateParameters(
	plan service.Plan,
	version string,
	dt *dbmsInstanceDetails,
	pp service.ProvisioningParameters,
) (map[string]interface{}, error) {

	td := plan.GetProperties().Extended["tierDetails"].(tierDetails)

	p := map[string]interface{}{}
	location := pp.GetString("location")
	p["location"] = location
	// Temporary workaround for Mooncake
	if location == "chinanorth" || location == "chinaeast" {
		p["family"] = "Gen4"
	} else {
		p["family"] = "Gen5"
	}
	p["sku"] = td.getSku(pp)
	p["tier"] = td.tierName
	p["cores"] = pp.GetInt64("cores")
	p["storage"] = pp.GetInt64("storage") * 1024 // storage is in MB to arm :/
	p["backupRetention"] = pp.GetInt64("backupRetention")
	if isGeoRedundentBackup(pp) {
		p["geoRedundantBackup"] = enabledARMString
	}
	p["version"] = version
	p["serverName"] = dt.ServerName
	p["administratorLogin"] = dt.AdministratorLogin
	p["administratorLoginPassword"] = string(dt.AdministratorLoginPassword)
	if isSSLRequired(pp) {
		p["sslEnforcement"] = enabledARMString
	} else {
		p["sslEnforcement"] = disabledARMString
	}
	firewallRulesParams := pp.GetObjectArray("firewallRules")
	firewallRules := make([]map[string]interface{}, len(firewallRulesParams))
	for i, firewallRuleParams := range firewallRulesParams {
		firewallRules[i] = firewallRuleParams.Data
	}
	p["firewallRules"] = firewallRules

	virtualNetworkRulesParams := pp.GetObjectArray("virtualNetworkRules")
	virtualNetworkRules := make([]map[string]interface{},
		len(virtualNetworkRulesParams))
	for i, virtualNetworkRulesParams := range virtualNetworkRulesParams {
		virtualNetworkRules[i] = virtualNetworkRulesParams.Data
	}
	p["virtualNetworkRules"] = virtualNetworkRules

	return p, nil
}
