package mysql

import (
	"context"
	"fmt"

	mysqlSDK "github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2017-12-01/mysql" // nolint: lll
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
		nameAvailability, err := checkNameAvailabilityClient.Execute(
			ctx,
			mysqlSDK.NameAvailabilityRequest{
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
