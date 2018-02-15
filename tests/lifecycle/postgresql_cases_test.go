// +build !unit

package lifecycle

import (
	"database/sql"
	"fmt"

	postgresSDK "github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-04-30-preview/postgresql" // nolint: lll
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/postgresqldb"
	_ "github.com/lib/pq" // Postgres SQL driver
)

func getPostgresqlCases(
	azureEnvironment azure.Environment,
	subscriptionID string,
	authorizer autorest.Authorizer,
	armDeployer arm.Deployer,
) ([]serviceLifecycleTestCase, error) {
	checkNameAvailabilityClient :=
		postgresSDK.NewCheckNameAvailabilityClientWithBaseURI(
			azureEnvironment.ResourceManagerEndpoint,
			subscriptionID,
		)
	checkNameAvailabilityClient.Authorizer = authorizer
	serversClient := postgresSDK.NewServersClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		subscriptionID,
	)
	serversClient.Authorizer = authorizer

	databasesClient := postgresSDK.NewDatabasesClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		subscriptionID,
	)
	databasesClient.Authorizer = authorizer

	module := postgresqldb.New(
		armDeployer,
		checkNameAvailabilityClient,
		serversClient,
		databasesClient,
	)

	return []serviceLifecycleTestCase{
		{
			module:          module,
			serviceID:       "b43b4bba-5741-4d98-a10b-17dc5cee0175",
			planID:          "b2ed210f-6a10-4593-a6c4-964e6b6fad62",
			description:     "all-in-one",
			location:        "southcentralus",
			testCredentials: testPostgreSQLCreds(),
			provisioningParameters: &postgresqldb.AllInOneProvisioningParameters{
				ServerProvisioningParameters: postgresqldb.ServerProvisioningParameters{ //nolint:lll
					FirewallIPStart: "0.0.0.0",
					FirewallIPEnd:   "255.255.255.255",
					SSLEnforcement:  "disabled",
				},
				DatabaseProvisioningParameters: postgresqldb.DatabaseProvisioningParameters{ //nolint:lll
					Extensions: []string{
						"uuid-ossp",
						"postgis",
					},
				},
			},
			bindingParameters: &postgresqldb.BindingParameters{},
		},
		{
			module:      module,
			serviceID:   "d3f74b44-79bc-4d1e-bf7d-c247c2b851f9",
			planID:      "bf389028-8dcc-433a-ab6f-0ee9b8db142f",
			description: "dbms-only",
			location:    "eastus",
			provisioningParameters: &postgresqldb.ServerProvisioningParameters{
				FirewallIPStart: "0.0.0.0",
				FirewallIPEnd:   "255.255.255.255",
				SSLEnforcement:  "disabled",
			},
			childTestCases: []*serviceLifecycleTestCase{
				{ // db only scenario
					module:            module,
					description:       "database-on-existing-server",
					serviceID:         "25434f16-d762-41c7-bbdd-8045d7f74ca6",
					planID:            "df6f5ef1-e602-406b-ba73-09c107d1e31b",
					location:          "", // This is actually irrelevant for this test
					bindingParameters: &postgresqldb.BindingParameters{},
					testCredentials:   testPostgreSQLCreds(),
					provisioningParameters: &postgresqldb.DatabaseProvisioningParameters{
						Extensions: []string{
							"uuid-ossp",
							"postgis",
						},
					},
				},
			},
		},
	}, nil
}

func testPostgreSQLCreds() func(credentials service.Credentials) error {
	return func(credentials service.Credentials) error {
		cdts, ok := credentials.(*postgresqldb.Credentials)
		if !ok {
			return fmt.Errorf(
				"error casting credentials as *postgresqldb.Credentials",
			)
		}
		connectionStrTemplate := "postgres://%s:%s@%s:%d/%s"
		connectionString := fmt.Sprintf(
			connectionStrTemplate,
			cdts.Username,
			cdts.Password,
			cdts.Host,
			cdts.Port,
			cdts.Database,
		)
		db, err := sql.Open("postgres", connectionString)

		if err != nil {
			return fmt.Errorf("error validating the database arguments: %s", err)
		}
		defer db.Close() // nolint: errcheck
		rows, err := db.Query(`
			SELECT * from pg_catalog.pg_tables
			WHERE
			schemaname != 'pg_catalog'
			AND schemaname != 'information_schema'
			`)
		if err != nil {
			return fmt.Errorf("error validating the database arguments: %s", err)
		}
		defer rows.Close() // nolint: errcheck
		if !rows.Next() {
			return fmt.Errorf(
				`error could not select from pg_catalog'`,
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
