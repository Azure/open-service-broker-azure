package mssql

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func buildDBMSGoTemplateParameters(
	dt *dbmsInstanceDetails,
	params service.ProvisioningParameters,
	version string,
) map[string]interface{} {
	p := map[string]interface{}{}
	p["serverName"] = dt.ServerName
	p["administratorLogin"] = dt.AdministratorLogin
	p["administratorLoginPassword"] = string(dt.AdministratorLoginPassword)
	p["version"] = version
	firewallRulesParams :=
		params.GetObjectArray("firewallRules")
	firewallRules := make([]map[string]interface{}, len(firewallRulesParams))
	for i, firewallRuleParams := range firewallRulesParams {
		firewallRules[i] = firewallRuleParams.Data
	}
	p["firewallRules"] = firewallRules
	return p
}

func buildDatabaseGoTemplateParameters(
	databaseName string,
	pp service.ProvisioningParameters,
	pd planDetails,
) (map[string]interface{}, error) {
	td, err := pd.getTierProvisionParameters(pp)
	if err != nil {
		return nil, err
	}
	p := map[string]interface{}{}
	p["databaseName"] = databaseName
	for key, value := range td {
		p[key] = value
	}
	return p, nil
}
