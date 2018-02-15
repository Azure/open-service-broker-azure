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
	"github.com/Azure/open-service-broker-azure/pkg/services/mysqldb"

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
	module := mysqldb.New(
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
			provisioningParameters: &mysqldb.ServerProvisioningParameters{
				SSLEnforcement:  "disabled",
				FirewallIPStart: "0.0.0.0",
				FirewallIPEnd:   "255.255.255.255",
			},
			bindingParameters: &mysqldb.BindingParameters{},
			testCredentials:   testMySQLCreds(),
		},
		{
			module:      module,
			description: "dbms server only",
			serviceID:   "30e7b836-199d-4335-b83d-adc7d23a95c2",
			planID:      "3f65ebf9-ac1d-4e77-b9bf-918889a4482b",
			location:    "eastus",
			provisioningParameters: &mysqldb.ServerProvisioningParameters{
				SSLEnforcement:  "disabled",
				FirewallIPStart: "0.0.0.0",
				FirewallIPEnd:   "255.255.255.255",
			},
			childTestCases: []*serviceLifecycleTestCase{
				{ // db only scenario
					module:                 module,
					description:            "database on new server",
					serviceID:              "6704ae59-3eae-49e9-82b4-4cbcc00edf08",
					planID:                 "ec77bd04-2107-408e-8fde-8100c1ce1f46",
					location:               "", // This is actually irrelevant for this test
					bindingParameters:      &mysqldb.BindingParameters{},
					testCredentials:        testMySQLCreds(),
					provisioningParameters: &mysqldb.DatabaseProvisioningParameters{},
				},
			},
		},
	}, nil
}

func testMySQLCreds() func(credentials service.Credentials) error {
	return func(credentials service.Credentials) error {

		cdts, ok := credentials.(*mysqldb.Credentials)
		if !ok {
			return fmt.Errorf("error casting credentials as *mssql.Credentials")
		}
		connectionTemplate := "%s:%s@tcp(%s:%d)/%s?allowNativePasswords=true"
		connectionString := fmt.Sprintf(
			connectionTemplate,
			cdts.Username,
			cdts.Password,
			cdts.Host,
			cdts.Port,
			cdts.Database,
		)
		db, err := sql.Open("mysql", connectionString)
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
