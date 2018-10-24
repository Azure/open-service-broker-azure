// +build !unit

package lifecycle

import (
	"context"
	"fmt"

	sqlSDK "github.com/Azure/azure-sdk-for-go/services/sql/mgmt/2017-03-01-preview/sql" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	_ "github.com/denisenkom/go-mssqldb" // MS SQL Driver
	uuid "github.com/satori/go.uuid"
)

var mssqlDBMSPairAlias = uuid.NewV4().String()

var mssqldrTestCases = []serviceLifecycleTestCase{
	{
		group:     "mssqldr",
		name:      "azure-sql-12-0-dr-dbms-pair-registered",
		serviceID: "00ce53a3-d6c3-4c24-8cb2-3f48d3b161d8",
		planID:    "5683ca92-372b-49a6-b7cd-96a14645ec15",
		preProvisionFns: []preProvisionFn{
			createSQLServerPair,
		},
		provisioningParameters: map[string]interface{}{
			"alias":             mssqlDBMSPairAlias,
			"primaryLocation":   "southcentralus",
			"secondaryLocation": "northcentralus",
		},
	},
}

func createSQLServerPair(
	ctx context.Context,
	resourceGroup string,
	_ *service.Instance,
	pp *map[string]interface{},
) error {
	azureConfig, err := getAzureConfig()
	if err != nil {
		return err
	}
	authorizer, err := getBearerTokenAuthorizer(azureConfig)
	if err != nil {
		return err
	}
	serversClient := sqlSDK.NewServersClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	serversClient.Authorizer = authorizer
	firewallRulesClient := sqlSDK.NewFirewallRulesClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	firewallRulesClient.Authorizer = authorizer

	priServerName := uuid.NewV4().String()
	priLocation := (*pp)["primaryLocation"].(string)
	priAdministratorLogin := generate.NewIdentifier()
	priAdministratorLoginPassword := generate.NewPassword()
	secServerName := uuid.NewV4().String()
	secLocation := (*pp)["secondaryLocation"].(string)
	secAdministratorLogin := generate.NewIdentifier()
	secAdministratorLoginPassword := generate.NewPassword()
	version := "12.0"
	(*pp)["primaryResourceGroup"] = resourceGroup
	(*pp)["primaryServer"] = priServerName
	(*pp)["primaryAdministratorLogin"] = priAdministratorLogin
	(*pp)["primaryAdministratorLoginPassword"] = priAdministratorLoginPassword
	(*pp)["secondaryResourceGroup"] = resourceGroup
	(*pp)["secondaryServer"] = secServerName
	(*pp)["secondaryAdministratorLogin"] = secAdministratorLogin
	(*pp)["secondaryAdministratorLoginPassword"] = secAdministratorLoginPassword

	priServer := sqlSDK.Server{
		Location: &priLocation,
		ServerProperties: &sqlSDK.ServerProperties{
			AdministratorLogin:         &priAdministratorLogin,
			AdministratorLoginPassword: &priAdministratorLoginPassword,
			Version:                    &version,
		},
	}
	priResult, err := serversClient.CreateOrUpdate(
		ctx,
		resourceGroup,
		priServerName,
		priServer,
	)
	if err != nil {
		return fmt.Errorf("error creating sql server: %s", err)
	}
	secServer := sqlSDK.Server{
		Location: &secLocation,
		ServerProperties: &sqlSDK.ServerProperties{
			AdministratorLogin:         &secAdministratorLogin,
			AdministratorLoginPassword: &secAdministratorLoginPassword,
			Version:                    &version,
		},
	}
	secResult, err := serversClient.CreateOrUpdate(
		ctx,
		resourceGroup,
		secServerName,
		secServer,
	)
	if err != nil {
		return fmt.Errorf("error creating sql server: %s", err)
	}
	if err := priResult.WaitForCompletion(ctx, serversClient.Client); err != nil {
		return fmt.Errorf("error creating sql server: %s", err)
	}
	if err := secResult.WaitForCompletion(ctx, serversClient.Client); err != nil {
		return fmt.Errorf("error creating sql server: %s", err)
	}
	startIPAddress := "0.0.0.0"
	endIPAddress := "255.255.255.255"
	firewallRule := sqlSDK.FirewallRule{
		FirewallRuleProperties: &sqlSDK.FirewallRuleProperties{
			StartIPAddress: &startIPAddress,
			EndIPAddress:   &endIPAddress,
		},
	}
	if _, err := firewallRulesClient.CreateOrUpdate(
		ctx,
		resourceGroup,
		priServerName,
		"all",
		firewallRule,
	); err != nil {
		return fmt.Errorf("error creating firewall rule: %s", err)
	}
	if _, err := firewallRulesClient.CreateOrUpdate(
		ctx,
		resourceGroup,
		secServerName,
		"all",
		firewallRule,
	); err != nil {
		return fmt.Errorf("error creating firewall rule: %s", err)
	}

	return nil
}
