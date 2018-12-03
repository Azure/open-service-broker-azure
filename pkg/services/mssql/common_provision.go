package mssql

import (
	"context"
	"fmt"

	sqlSDK "github.com/Azure/azure-sdk-for-go/services/sql/mgmt/2017-03-01-preview/sql" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func buildDBMSGoTemplateParameters(
	dt *dbmsInstanceDetails,
	params service.ProvisioningParameters,
	version string,
) (map[string]interface{}, error) {
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

	return p, nil
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

func setConnectionPolicy(
	ctx context.Context,
	client *sqlSDK.ServerConnectionPoliciesClient,
	resourceGroupName string,
	serverName string,
	location string,
	connectionPolicy string,
) error {
	var serverConnectionPolicy sqlSDK.ServerConnectionPolicy
	serverConnectionPolicy.Location = &location
	switch connectionPolicy {
	case string(sqlSDK.ServerConnectionTypeDefault):
		serverConnectionPolicy.ServerConnectionPolicyProperties =
			&sqlSDK.ServerConnectionPolicyProperties{
				ConnectionType: sqlSDK.ServerConnectionTypeDefault,
			}
	case string(sqlSDK.ServerConnectionTypeProxy):
		serverConnectionPolicy.ServerConnectionPolicyProperties =
			&sqlSDK.ServerConnectionPolicyProperties{
				ConnectionType: sqlSDK.ServerConnectionTypeProxy,
			}
	case string(sqlSDK.ServerConnectionTypeRedirect):
		serverConnectionPolicy.ServerConnectionPolicyProperties =
			&sqlSDK.ServerConnectionPolicyProperties{
				ConnectionType: sqlSDK.ServerConnectionTypeRedirect,
			}
	default:
		return fmt.Errorf("no such connection policy")
	}
	_, err := client.CreateOrUpdate(
		ctx,
		resourceGroupName,
		serverName,
		serverConnectionPolicy,
	)
	return err
}
