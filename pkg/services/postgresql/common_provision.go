package postgresql

import (
	"bytes"
	"context"
	"fmt"
	"net"

	postgresSDK "github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-04-30-preview/postgresql" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/slice"
	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq" // Postgres SQL driver
	uuid "github.com/satori/go.uuid"
)

const (
	enabledParamString  = "enabled"
	disabledParamString = "disabled"
	enabledARMString    = "Enabled"
	disabledARMString   = "Disabled"
)

func getAvailableServerName(
	ctx context.Context,
	checkNameAvailabilityClient postgresSDK.CheckNameAvailabilityClient,
) (string, error) {
	for {
		serverName := uuid.NewV4().String()
		nameAvailability, err := checkNameAvailabilityClient.Execute(
			ctx,
			postgresSDK.NameAvailabilityRequest{
				Name: &serverName,
			},
		)
		if err != nil {
			return "", fmt.Errorf(
				"error determining server name availability: %s",
				err,
			)
		}
		if *nameAvailability.NameAvailable {
			return serverName, nil
		}
	}
}

func validateDBMSProvisionParameters(
	plan service.Plan,
	pp dbmsProvisioningParameters,
) error {
	if plan == nil {
		return fmt.Errorf("plan invalid")
	}
	s, ok := plan.GetProperties().Extended["provisionSchema"]
	if !ok {
		return fmt.Errorf("invalid plan, schema not found")
	}
	schema := s.(planSchema)
	if pp.SSLEnforcement != "" &&
		!slice.ContainsString(schema.allowedSSLEnforcement, pp.SSLEnforcement) {
		return service.NewValidationError(
			"sslEnforcement",
			fmt.Sprintf(`invalid option: "%s"`, pp.SSLEnforcement),
		)
	}
	for _, firewallRule := range pp.FirewallRules {
		if firewallRule.Name == "" {
			return service.NewValidationError(
				"ruleName",
				"must be set",
			)
		}
		if firewallRule.StartIP != "" || firewallRule.EndIP != "" {
			if firewallRule.StartIP == "" {
				return service.NewValidationError(
					"startIPAddress",
					"must be set when endIPAddress is set",
				)
			}
			if firewallRule.EndIP == "" {
				return service.NewValidationError(
					"endIPAddress",
					"must be set when startIPAddress is set",
				)
			}
		}
		startIP := net.ParseIP(firewallRule.StartIP)
		if firewallRule.StartIP != "" && startIP == nil {
			return service.NewValidationError(
				"startIPAddress",
				fmt.Sprintf(`invalid value: "%s"`, firewallRule.StartIP),
			)
		}
		endIP := net.ParseIP(firewallRule.StartIP)
		if firewallRule.EndIP != "" && endIP == nil {
			return service.NewValidationError(
				"endIPAddress",
				fmt.Sprintf(
					`invalid value: "%s"`,
					firewallRule.EndIP,
				),
			)
		}
		// The net.IP.To4 method returns a 4 byte representation of an IPv4 address.
		// Once converted,comparing two IP addresses can be done by using the
		// bytes. Compare function. Per the ARM template documentation,
		// startIP must be <= endIP.
		startBytes := startIP.To4()
		endBytes := endIP.To4()
		if bytes.Compare(startBytes, endBytes) > 0 {
			return service.NewValidationError(
				"endIPAddress",
				fmt.Sprintf(`invalid value: "%s". must be 
				greater than or equal to startIPAddress`,
					firewallRule.EndIP),
			)
		}
	}

	// hardware family
	if pp.HardwareFamily != "" &&
		!slice.ContainsString(schema.allowedHardware, pp.HardwareFamily) {
		return service.NewValidationError(
			"hardwareFamily",
			fmt.Sprintf(`invalid value: "%s"`, pp.HardwareFamily))
	}

	// cores
	if pp.Cores != nil && !slice.ContainsInt64(schema.allowedCores, *pp.Cores) {
		return service.NewValidationError(
			"cores",
			fmt.Sprintf(`invalid value: "%d"`, *pp.Cores),
		)
	}

	// storage
	if pp.Storage != nil &&
		(*pp.Storage < schema.minStorage || *pp.Storage > schema.maxStorage) {
		return service.NewValidationError(
			"storage",
			fmt.Sprintf(`invalid value: "%d"`, *pp.Storage),
		)
	}

	// backupRetation
	if pp.BackupRetention != nil &&
		(*pp.BackupRetention < schema.minBackupRetention ||
			*pp.BackupRetention > schema.maxBackupRetention) {
		return service.NewValidationError(
			"backupRetention",
			fmt.Sprintf(`invalid value: "%d"`, *pp.BackupRetention),
		)
	}

	// backupRedundancy
	if pp.BackupRedundancy != "" &&
		!slice.ContainsString(schema.allowedBackupRedundancy, pp.BackupRedundancy) {
		return service.NewValidationError(
			"backupRedundancy",
			fmt.Sprintf(`invalid value: "%s"`, pp.BackupRedundancy))
	}

	return nil
}

func setupDatabase(
	enforceSSL bool,
	serverName string,
	administratorLoginPassword string,
	fullyQualifiedDomainName string,
	dbName string,
) error {
	db, err := getDBConnection(
		enforceSSL,
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
		fmt.Sprintf("create role %s", dbName),
	); err != nil {
		return fmt.Errorf(`error creating role "%s": %s`, dbName, err)
	}
	if _, err = tx.Exec(
		fmt.Sprintf("grant %s to postgres", dbName),
	); err != nil {
		return fmt.Errorf(
			`error adding role "%s" to role "postgres": %s`,
			dbName,
			err,
		)
	}
	if _, err = tx.Exec(
		fmt.Sprintf(
			"alter database %s owner to %s",
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
	serverName string,
	administratorLoginPassword string,
	fullyQualifiedDomainName string,
	dbName string,
	extensions []string,
) error {
	db, err := getDBConnection(
		enforceSSL,
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
	instance service.Instance,
) (map[string]interface{}, error) {

	plan := instance.Plan
	dt := dbmsInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, err
	}
	sdt := secureAllInOneInstanceDetails{}
	if err := service.GetStructFromMap(instance.SecureDetails, &sdt); err != nil {
		return nil, err
	}
	pp := dbmsProvisioningParameters{}
	if err :=
		service.GetStructFromMap(instance.ProvisioningParameters, &pp); err != nil {
		return nil, err
	}

	schema := plan.GetProperties().Extended["provisionSchema"].(planSchema)

	p := map[string]interface{}{}
	p["sku"] = schema.getSku(pp)
	p["tier"] = plan.GetProperties().Extended["tier"]
	p["cores"] = schema.getCores(pp)
	p["storage"] = schema.getStorage(pp) * 1024 //storage is in MB to arm :/
	p["backupRetention"] = schema.getBackupRetention(pp)
	p["hardwareFamily"] = schema.getHardwareFamily(pp)
	if schema.isGeoRedundentBackup(pp) {
		p["geoRedundantBackup"] = enabledARMString
	}
	p["version"] = instance.Service.GetProperties().Extended["version"]
	p["serverName"] = dt.ServerName
	p["administratorLoginPassword"] = sdt.AdministratorLoginPassword
	if schema.isSSLRequired(pp) {
		p["sslEnforcement"] = enabledARMString
	} else {
		p["sslEnforcement"] = disabledARMString
	}
	p["firewallRules"] = schema.getFirewallRules(pp)

	return p, nil
}
