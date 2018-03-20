// +build !unit

package lifecycle

import (
	"database/sql"
	"fmt"

	mysqlSDK "github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2017-04-30-preview/mysql" // nolint: lll
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	mysql "github.com/Azure/open-service-broker-azure/pkg/services/mysql"

	_ "github.com/go-sql-driver/mysql"
)

func getMysqlCases(
	azureEnvironment azure.Environment,
	subscriptionID string,
	authorizer autorest.Authorizer,
	armDeployer arm.Deployer,
) ([]serviceLifecycleTestCase, error) {
	checkNameAvailabilityClient :=
		mysqlSDK.NewCheckNameAvailabilityClientWithBaseURI(
			azureEnvironment.ResourceManagerEndpoint,
			subscriptionID,
		)
	checkNameAvailabilityClient.Authorizer = authorizer
	serversClient := mysqlSDK.NewServersClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		subscriptionID,
	)
	serversClient.Authorizer = authorizer
	databasesClient := mysqlSDK.NewDatabasesClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		subscriptionID,
	)
	databasesClient.Authorizer = authorizer
	module := mysql.New(
		azureEnvironment,
		armDeployer,
		checkNameAvailabilityClient,
		serversClient,
		databasesClient,
	)
	return []serviceLifecycleTestCase{
		{
			module:      module,
			description: "server and database (all-in-one)",
			serviceID:   "997b8372-8dac-40ac-ae65-758b4a5075a5",
			planID:      "427559f1-bf2a-45d3-8844-32374a3e58aa",
			location:    "southcentralus",
			provisioningParameters: service.CombinedProvisioningParameters{
				"sslEnforcement": "disabled",
				"firewallRules": []map[string]string{
					{
						"name":           "AllowSome",
						"startIPAddress": "0.0.0.0",
						"endIPAddress":   "35.0.0.0",
					},
					{
						"name":           "AllowMore",
						"startIPAddress": "35.0.0.1",
						"endIPAddress":   "255.255.255.255",
					},
				},
			},
			bindingParameters: nil,
			testCredentials:   testMySQLCreds(),
		},
		{
			module:      module,
			description: "dbms server only",
			serviceID:   "30e7b836-199d-4335-b83d-adc7d23a95c2",
			planID:      "3f65ebf9-ac1d-4e77-b9bf-918889a4482b",
			location:    "eastus",
			provisioningParameters: service.CombinedProvisioningParameters{
				"firewallRules": []map[string]string{
					{
						"name":           "AllowAll",
						"startIPAddress": "0.0.0.0",
						"endIPAddress":   "255.255.255.255",
					},
				},
			},
			childTestCases: []*serviceLifecycleTestCase{
				{ // database only scenario
					module:                 module,
					description:            "database on new server",
					serviceID:              "6704ae59-3eae-49e9-82b4-4cbcc00edf08",
					planID:                 "ec77bd04-2107-408e-8fde-8100c1ce1f46",
					location:               "", // This is actually irrelevant for this test
					bindingParameters:      nil,
					testCredentials:        testMySQLCreds(),
					provisioningParameters: nil,
				},
			},
		},
	}, nil
}

func testMySQLCreds() func(credentials map[string]interface{}) error {
	return func(credentials map[string]interface{}) error {

		var connectionStrTemplate string
		if credentials["sslRequired"].(bool) {
			connectionStrTemplate =
				"%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true&tls=true"
		} else {
			connectionStrTemplate =
				"%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true"
		}

		db, err := sql.Open("mysql", fmt.Sprintf(
			connectionStrTemplate,
			credentials["username"].(string),
			credentials["password"].(string),
			credentials["host"].(string),
			credentials["database"].(string),
		))
		if err != nil {
			return fmt.Errorf("error validating the database arguments: %s", err)
		}
		defer db.Close() // nolint: errcheck
		rows, err := db.Query("SELECT * from INFORMATION_SCHEMA.TABLES")
		if err != nil {
			return fmt.Errorf("error validating the database arguments: %s", err)
		}
		defer rows.Close() // nolint: errcheck
		if !rows.Next() {
			return fmt.Errorf(
				`error could not select from INFORMATION_SCHEMA.TABLES'`,
			)
		}
		if err := rows.Err(); err != nil {
			return fmt.Errorf(
				`error iterating rows`,
			)
		}
		return nil
	}
}
