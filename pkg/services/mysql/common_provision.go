package mysql

import (
	"context"
	"fmt"

	mysqlSDK "github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2017-12-01/mysql" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
	uuid "github.com/satori/go.uuid"
)

const (
	enabledARMString  = "Enabled"
	disabledARMString = "Disabled"
)

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

	return p, nil
}

func getAvailableServerName(
	ctx context.Context,
	checkNameAvailabilityClient mysqlSDK.CheckNameAvailabilityClient,
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

// generateDBMSInstanceDetail will read information
// from instance provision parameters, and generate
// a dbmsInstanceDetail. This method is expected to
// be invoked by preProvision step of all-in-one and
// dbms.
func generateDBMSInstanceDetails(
	ctx context.Context,
	instance service.Instance,
	checkNameAvailabilityClient mysqlSDK.CheckNameAvailabilityClient,
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
			return nil, fmt.Errorf("server name %s is not available", serverName)
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
	// use it; else, use default value "azureuser".
	adminAccountSettings := pp.GetObject("adminAccountSettings")
	adminUsername := adminAccountSettings.GetString("adminUsername")
	if adminUsername == "" {
		adminUsername = "azureuser"
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

func isServerNameAvailable(
	ctx context.Context,
	serverName string,
	checkNameAvailabilityClient mysqlSDK.CheckNameAvailabilityClient,
) (bool, error) {
	nameAvailability, err := checkNameAvailabilityClient.Execute(
		ctx,
		mysqlSDK.NameAvailabilityRequest{
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
