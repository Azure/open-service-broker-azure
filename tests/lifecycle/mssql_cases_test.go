// +build !unit

package lifecycle

import (
	"database/sql"
	"fmt"
	"net/url"

	sqlSDK "github.com/Azure/azure-sdk-for-go/services/sql/mgmt/2017-03-01-preview/sql" // nolint: lll
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/mssql"
	_ "github.com/denisenkom/go-mssqldb" // MS SQL Driver
)

func getMssqlCases(
	azureEnvironment azure.Environment,
	subscriptionID string,
	authorizer autorest.Authorizer,
	armDeployer arm.Deployer,
) ([]serviceLifecycleTestCase, error) {
	sqlServersClient := sqlSDK.NewServersClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		subscriptionID,
	)
	sqlServersClient.Authorizer = authorizer
	sqlDatabasesClient := sqlSDK.NewDatabasesClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		subscriptionID,
	)
	sqlDatabasesClient.Authorizer = authorizer
	module := mssql.New(
		azureEnvironment,
		armDeployer,
		sqlServersClient,
		sqlDatabasesClient,
	)
	return []serviceLifecycleTestCase{
		{ // all-in-one scenario
			module:      module,
			description: "new server and database (all in one)",
			serviceID:   "fb9bc99e-0aa9-11e6-8a8a-000d3a002ed5",
			planID:      "3819fdfa-0aaa-11e6-86f4-000d3a002ed5",
			location:    "southcentralus",
			provisioningParameters: service.CombinedProvisioningParameters{
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
			testCredentials:   testMsSQLCreds(),
		},
		{ // server only scenario
			module:      module,
			description: "new server with database child test",
			serviceID:   "a7454e0e-be2c-46ac-b55f-8c4278117525",
			planID:      "24f0f42e-1ab3-474e-a5ca-b943b2c48eee",
			location:    "southcentralus",
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
				{ // db only scenario
					module:            module,
					description:       "database on new server",
					serviceID:         "2bbc160c-e279-4757-a6b6-4c0a4822d0aa",
					planID:            "8fa8d759-c142-45dd-ae38-b93482ddc04a",
					location:          "", // This is actually irrelevant for this test
					bindingParameters: nil,
					testCredentials:   testMsSQLCreds(),
				},
			},
		},
	}, nil
}

func testMsSQLCreds() func(credentials map[string]interface{}) error {
	return func(credentials map[string]interface{}) error {
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
}
