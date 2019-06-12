// +build !unit

package lifecycle

import (
	"context"
	"fmt"
	"time"

	sqlSDK "github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql" // nolint: lll
	"github.com/Azure/go-autorest/autorest/to"
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
			"secondaryLocation": "eastus",
		},
		childTestCases: []*serviceLifecycleTestCase{
			{
				group:           "mssqldr",
				name:            "azure-sql-12-0-dr-database-pair",
				serviceID:       "2eb94a7e-5a7c-46f9-b9d2-ff769f215845",
				planID:          "edce3e74-69eb-4524-aabb-f2c4a7ee9398",
				testCredentials: testMsSQLCreds,
				provisioningParameters: map[string]interface{}{
					"parentAlias":   mssqlDBMSPairAlias,
					"failoverGroup": uuid.NewV4().String(),
					"database":      uuid.NewV4().String(),
				},
				updatingParameters: map[string]interface{}{
					"dtus": 50,
				},
			},
			{
				group:           "mssqldr",
				name:            "azure-sql-12-0-dr-database-pair-registered",
				serviceID:       "8480271a-f4c7-4232-b2b7-7f33391728f7",
				planID:          "9e05f8b7-27ce-4fb4-b889-e7b2f8575df7",
				testCredentials: testMsSQLCreds,
				preProvisionFns: []preProvisionFn{
					createSQLDatabasePair,
				},
				provisioningParameters: map[string]interface{}{
					"parentAlias": mssqlDBMSPairAlias,
					// "failoverGroup" and "database" are specified
					// by createSQLDatabasePair
				},
			},
			{
				group:           "mssqldr",
				name:            "azure-sql-12-0-dr-database-pair-from-existing",
				serviceID:       "e18a9861-5740-4e1a-9bd0-6f0fc3e4d12f",
				planID:          "af66e3e3-c500-4042-879e-5a6d47901d1c",
				testCredentials: testMsSQLCreds,
				preProvisionFns: []preProvisionFn{
					createSQLDatabasePair,
				},
				provisioningParameters: map[string]interface{}{
					"parentAlias": mssqlDBMSPairAlias,
					// "failoverGroup" and "database" are specified
					// by createSQLDatabasePair
				},
				updatingParameters: map[string]interface{}{
					"dtus": 50,
				},
			},
			{
				group:           "mssqldr",
				name:            "azure-sql-12-0-dr-database-pair-from-existing-primary", // nolint: lll
				serviceID:       "505ae87a-5cd8-4aeb-b7ea-809dd249dc1f",
				planID:          "8ec86bea-42f6-4805-b3e9-506eaebbf9e0",
				testCredentials: testMsSQLCreds,
				preProvisionFns: []preProvisionFn{
					createPrimarySQLDatabase,
				},
				provisioningParameters: map[string]interface{}{
					"parentAlias":   mssqlDBMSPairAlias,
					"failoverGroup": uuid.NewV4().String(),
					// "database" is specified by createPrimarySQLDatabase
				},
				updatingParameters: map[string]interface{}{
					"dtus": 50,
				},
			},
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
	if err := priResult.WaitForCompletionRef(ctx, serversClient.Client); err != nil {
		return fmt.Errorf("error creating sql server: %s", err)
	}
	if err := secResult.WaitForCompletionRef(ctx, serversClient.Client); err != nil {
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

func createSQLDatabasePair(
	ctx context.Context,
	resourceGroup string,
	parent *service.Instance,
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
	databasesClient := sqlSDK.NewDatabasesClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	databasesClient.Authorizer = authorizer
	failoverGroupsClient := sqlSDK.NewFailoverGroupsClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	failoverGroupsClient.Authorizer = authorizer

	dtMap, err := service.GetMapFromStruct(parent.Details)
	if err != nil {
		return err
	}
	priServerName := dtMap["primaryServer"].(string)
	secServerName := dtMap["secondaryServer"].(string)
	priLocation := parent.ProvisioningParameters.GetString("primaryLocation")   // nolint: lll
	secLocation := parent.ProvisioningParameters.GetString("secondaryLocation") // nolint: lll
	databaseName := generate.NewIdentifier()
	database := sqlSDK.Database{
		Location: &priLocation,
	}
	failoverGroupName := generate.NewIdentifier()
	failoverGroup := sqlSDK.FailoverGroup{
		FailoverGroupProperties: &sqlSDK.FailoverGroupProperties{
			ReadWriteEndpoint: &sqlSDK.FailoverGroupReadWriteEndpoint{
				FailoverPolicy:                         sqlSDK.Automatic,
				FailoverWithDataLossGracePeriodMinutes: to.Int32Ptr(60),
			},
			ReadOnlyEndpoint: &sqlSDK.FailoverGroupReadOnlyEndpoint{
				FailoverPolicy: sqlSDK.ReadOnlyEndpointFailoverPolicyDisabled,
			},
			ReplicationRole: sqlSDK.Primary,
			PartnerServers: &[]sqlSDK.PartnerInfo{
				{
					ID: to.StringPtr(fmt.Sprintf(
						"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql"+
							"/servers/%s",
						azureConfig.SubscriptionID,
						resourceGroup,
						secServerName,
					)),
				},
			},
			Databases: &[]string{
				fmt.Sprintf(
					"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Sql"+
						"/servers/%s/databases/%s",
					azureConfig.SubscriptionID,
					resourceGroup,
					priServerName,
					databaseName,
				),
			},
		},
	}
	(*pp)["database"] = databaseName
	(*pp)["failoverGroup"] = failoverGroupName

	priResult, err := databasesClient.CreateOrUpdate(
		ctx,
		resourceGroup,
		priServerName,
		databaseName,
		database,
	)
	if err != nil {
		return fmt.Errorf("error creating sql database: %s", err)
	}
	if err = priResult.WaitForCompletionRef(
		ctx,
		databasesClient.Client,
	); err != nil {
		return fmt.Errorf("error creating sql database: %s", err)
	}
	fgResult, err := failoverGroupsClient.CreateOrUpdate(
		ctx,
		resourceGroup,
		priServerName,
		failoverGroupName,
		failoverGroup,
	)
	if err != nil {
		return fmt.Errorf("error creating sql failover group: %s", err)
	}
	if err = fgResult.WaitForCompletionRef(
		ctx,
		failoverGroupsClient.Client,
	); err != nil {
		return fmt.Errorf("error creating sql failover group: %s", err)
	}
	// The secondary db is created by the failover group creation. But in a very
	//   short time after creating failover group, you can't get the secondary db.
	//   Here is the workaround to ensure that the getting of the secondary
	//   db could succeed.
	time.Sleep(time.Second * 10)
	database = sqlSDK.Database{
		Location: &secLocation,
	}
	secResult, err := databasesClient.CreateOrUpdate(
		ctx,
		resourceGroup,
		secServerName,
		databaseName,
		database,
	)
	if err != nil {
		return fmt.Errorf("error creating sql database: %s", err)
	}
	if err := secResult.WaitForCompletionRef(
		ctx,
		databasesClient.Client,
	); err != nil {
		return fmt.Errorf("error creating sql database: %s", err)
	}
	return nil
}

func createPrimarySQLDatabase(
	ctx context.Context,
	resourceGroup string,
	parent *service.Instance,
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
	databasesClient := sqlSDK.NewDatabasesClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	databasesClient.Authorizer = authorizer

	dtMap, err := service.GetMapFromStruct(parent.Details)
	if err != nil {
		return err
	}
	serverName := dtMap["primaryServer"].(string)
	databaseName := generate.NewIdentifier()
	location := parent.ProvisioningParameters.GetString("primaryLocation")
	database := sqlSDK.Database{
		Location: &location,
	}
	(*pp)["database"] = databaseName

	result, err := databasesClient.CreateOrUpdate(
		ctx,
		resourceGroup,
		serverName,
		databaseName,
		database,
	)
	if err != nil {
		return fmt.Errorf("error creating sql database: %s", err)
	}
	if err := result.WaitForCompletionRef(
		ctx,
		databasesClient.Client,
	); err != nil {
		return fmt.Errorf("error creating sql database: %s", err)
	}
	return nil
}
