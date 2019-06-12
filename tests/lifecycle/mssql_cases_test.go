// +build !unit

package lifecycle

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	sqlSDK "github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2017-03-01-preview/sql" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	_ "github.com/denisenkom/go-mssqldb" // MS SQL Driver
	uuid "github.com/satori/go.uuid"
)

var mssqlDBMSAlias = uuid.NewV4().String()
var mssqlDBMSRegisteredAlias = uuid.NewV4().String()

var mssqlTestCases = []serviceLifecycleTestCase{
	{ // all-in-one scenario (dtu-based)
		group:     "mssql",
		name:      "all-in-one (DTU)",
		serviceID: "fb9bc99e-0aa9-11e6-8a8a-000d3a002ed5",
		planID:    "2497b7f3-341b-4ac6-82fb-d4a48c005e19",
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"dtus":     200,
			"firewallRules": []interface{}{
				map[string]interface{}{
					"name":           "AllowSome",
					"startIPAddress": "0.0.0.0",
					"endIPAddress":   "35.0.0.0",
				},
				map[string]interface{}{
					"name":           "AllowMore",
					"startIPAddress": "35.0.0.1",
					"endIPAddress":   "255.255.255.255",
				},
			},
			"connectionPolicy": "Proxy",
		},
		updatingParameters: map[string]interface{}{
			"dtus": 400,
			"firewallRules": []interface{}{
				map[string]interface{}{
					"name":           "AllowSome2",
					"startIPAddress": "0.0.0.0",
					"endIPAddress":   "45.0.0.0",
				},
				map[string]interface{}{
					"name":           "AllowMore2",
					"startIPAddress": "45.0.0.1",
					"endIPAddress":   "255.255.255.255",
				},
			},
			"connectionPolicy": "Default",
		},
		testCredentials: testMsSQLCreds,
	},
	{ // all-in-one scenario (vcore-based)
		group:     "mssql",
		name:      "all-in-one (vCore)",
		serviceID: "fb9bc99e-0aa9-11e6-8a8a-000d3a002ed5",
		planID:    "c77e86af-f050-4457-a2ff-2b48451888f3",
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"cores":    4,
			"storage":  25,
			"firewallRules": []interface{}{
				map[string]interface{}{
					"name":           "AllowAll",
					"startIPAddress": "0.0.0.0",
					"endIPAddress":   "255.255.255.255",
				},
			},
		},
		updatingParameters: map[string]interface{}{
			"cores":   8,
			"storage": 50,
			"firewallRules": []interface{}{
				map[string]interface{}{
					"name":           "AllowSome2",
					"startIPAddress": "0.0.0.0",
					"endIPAddress":   "45.0.0.0",
				},
				map[string]interface{}{
					"name":           "AllowMore2",
					"startIPAddress": "45.0.0.1",
					"endIPAddress":   "255.255.255.255",
				},
			},
			"connectionPolicy": "Redirect",
		},
	},
	{ // dbms only scenario
		group:     "mssql",
		name:      "dbms-only",
		serviceID: "a7454e0e-be2c-46ac-b55f-8c4278117525",
		planID:    "24f0f42e-1ab3-474e-a5ca-b943b2c48eee",
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"alias":    mssqlDBMSAlias,
			"firewallRules": []interface{}{
				map[string]interface{}{
					"name":           "AllowAll",
					"startIPAddress": "0.0.0.0",
					"endIPAddress":   "255.255.255.255",
				},
			},
		},
		updatingParameters: map[string]interface{}{
			"connectionPolicy": "Proxy",
		},
		childTestCases: []*serviceLifecycleTestCase{
			{ // db only scenario (dtu-based)
				group:           "mssql",
				name:            "database-only (DTU)",
				serviceID:       "2bbc160c-e279-4757-a6b6-4c0a4822d0aa",
				planID:          "9d36b6b3-b5f3-4907-a713-5cc13b785409",
				testCredentials: testMsSQLCreds,
				provisioningParameters: map[string]interface{}{
					"parentAlias": mssqlDBMSAlias,
				},
				updatingParameters: map[string]interface{}{
					"dtus": 50,
				},
			},
			{ // db only scenario (vcore-based)
				group:           "mssql",
				name:            "database-only (vCore)",
				serviceID:       "2bbc160c-e279-4757-a6b6-4c0a4822d0aa",
				planID:          "da591616-77a1-4df8-a493-6c119649bc6b",
				testCredentials: testMsSQLCreds,
				provisioningParameters: map[string]interface{}{
					"parentAlias": mssqlDBMSAlias,
					"cores":       2,
					"storage":     10,
				},
				updatingParameters: map[string]interface{}{
					"cores":   4,
					"storage": 20,
				},
			},
			{ // db only from existing scenario (dtu-based)
				group:           "mssql",
				name:            "database-only-fe (DTU)",
				serviceID:       "b0b2a2f7-9b5e-4692-8b94-24fe2f6a9a8e",
				planID:          "e5804586-625a-4f67-996f-ca19a14711cc",
				testCredentials: testMsSQLCreds,
				preProvisionFns: []preProvisionFn{
					createSQLDatabase,
				},
				provisioningParameters: map[string]interface{}{
					"parentAlias": mssqlDBMSAlias,
				},
			},
		},
	},
	{ // dbms only registered scenario
		group:     "mssql",
		name:      "dbms-only-registered",
		serviceID: "c9bd94e7-5b7d-4b20-96e6-c5678f99d997",
		planID:    "4e95e962-0142-4117-b212-bcc7aec7f6c2",
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"alias":    mssqlDBMSRegisteredAlias,
		},
		preProvisionFns: []preProvisionFn{
			createSQLServer,
		},
		childTestCases: []*serviceLifecycleTestCase{
			{ // dtu db only scenario
				group:           "mssql",
				name:            "database-only (DTU)",
				serviceID:       "2bbc160c-e279-4757-a6b6-4c0a4822d0aa",
				planID:          "8fa8d759-c142-45dd-ae38-b93482ddc04a",
				testCredentials: testMsSQLCreds,
				provisioningParameters: map[string]interface{}{
					"parentAlias": mssqlDBMSRegisteredAlias,
				},
			},
			{ // vcore db only scenario
				group:           "mssql",
				name:            "database-only (vCore)",
				serviceID:       "2bbc160c-e279-4757-a6b6-4c0a4822d0aa",
				planID:          "da591616-77a1-4df8-a493-6c119649bc6b",
				testCredentials: testMsSQLCreds,
				provisioningParameters: map[string]interface{}{
					"parentAlias": mssqlDBMSRegisteredAlias,
					"cores":       int64(2),
					"storage":     int64(10),
				},
			},
		},
	},
}

func createSQLServer(
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

	serverName := uuid.NewV4().String()
	administratorLogin := generate.NewIdentifier()
	administratorLoginPassword := generate.NewPassword()
	version := "12.0"
	location := (*pp)["location"].(string)
	(*pp)["server"] = serverName
	(*pp)["administratorLogin"] = administratorLogin
	(*pp)["administratorLoginPassword"] = administratorLoginPassword

	server := sqlSDK.Server{
		Location: &location,
		ServerProperties: &sqlSDK.ServerProperties{
			AdministratorLogin:         &administratorLogin,
			AdministratorLoginPassword: &administratorLoginPassword,
			Version:                    &version,
		},
	}
	result, err := serversClient.CreateOrUpdate(
		ctx,
		resourceGroup,
		serverName,
		server,
	)
	if err != nil {
		return fmt.Errorf("error creating sql server: %s", err)
	}
	if err := result.WaitForCompletionRef(ctx, serversClient.Client); err != nil {
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
		serverName,
		"all",
		firewallRule,
	); err != nil {
		return fmt.Errorf("error creating firewall rule: %s", err)
	}
	return nil
}

func createSQLDatabase(
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

	pdtMap, err := service.GetMapFromStruct(parent.Details)
	if err != nil {
		return err
	}
	serverName := pdtMap["server"].(string)
	databaseName := generate.NewIdentifier()
	location := parent.ProvisioningParameters.GetString("location")
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
	if err := result.WaitForCompletionRef(ctx, databasesClient.Client); err != nil {
		return fmt.Errorf("error creating sql database: %s", err)
	}
	return nil
}

func testMsSQLCreds(credentials map[string]interface{}) error {
	query := url.Values{}
	query.Add("database", credentials["database"].(string))
	query.Add("encrypt", "true")
	query.Add("TrustServerCertificate", "true")

	u := &url.URL{
		Scheme: "sqlserver",
		User: url.UserPassword(
			credentials["username"].(string),
			credentials["password"].(string),
		),
		Host: fmt.Sprintf(
			"%s:%d",
			credentials["host"].(string),
			int(credentials["port"].(float64)),
		),
		RawQuery: query.Encode(),
	}

	db, err := sql.Open("mssql", u.String())
	if err != nil {
		return fmt.Errorf("error validating the database arguments: %s", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("error connecting to the database: %s", err)
	}
	defer db.Close() // nolint: errcheck

	rows, err := db.Query("SELECT 1 FROM fn_my_permissions (NULL, 'DATABASE') WHERE permission_name='CONTROL'") // nolint: lll
	if err != nil {
		return fmt.Errorf(
			`error querying SELECT from table fn_my_permissions: %s`,
			err,
		)
	}
	defer rows.Close() // nolint: errcheck
	if !rows.Next() {
		return fmt.Errorf(
			`error user doesn't have permission 'CONTROL'`,
		)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf(
			`error iterating rows`,
		)
	}

	return nil
}
