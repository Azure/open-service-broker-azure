// +build !unit

package lifecycle

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	ss "github.com/Azure/open-service-broker-azure/pkg/azure/mssql"
	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/mssqldb"
	_ "github.com/denisenkom/go-mssqldb" // MS SQL Driver
	uuid "github.com/satori/go.uuid"
)

// nolint: lll
var armTemplateNewServerBytes = []byte(`
{
	"$schema": "http://schema.management.azure.com/schemas/2014-04-01-preview/deploymentTemplate.json#",
	"contentVersion": "1.0.0.0",
	"parameters": {
		"location": {
			"type": "string"
		},
		"serverName": {
			"type": "string"
		},
		"administratorLogin": {
			"type": "string"
		},
		"administratorLoginPassword": {
			"type": "securestring"
		},
		"tags": {
			"type": "object"
		}
	},
	"variables": {
		"SQLapiVersion": "2014-04-01"
	},
	"resources": [
		{
			"type": "Microsoft.Sql/servers",
			"name": "[parameters('serverName')]",
			"apiVersion": "[variables('SQLapiVersion')]",
			"location": "[parameters('location')]",
			"properties": {
				"administratorLogin": "[parameters('administratorLogin')]",
				"administratorLoginPassword": "[parameters('administratorLoginPassword')]",
				"version": "12.0"
			},
			"resources": [
				{
					"type": "firewallrules",
					"name": "all",
					"apiVersion": "[variables('SQLapiVersion')]",
					"location": "[parameters('location')]",
					"properties": {
						"startIpAddress": "0.0.0.0",
						"endIpAddress": "255.255.255.255"
					},
					"dependsOn": [
						"[concat('Microsoft.Sql/servers/', parameters('serverName'))]"
					]
				}
			]
		}
	],
	"outputs": {
	}
}
`)

func getMssqlCases(
	armDeployer arm.Deployer,
	resourceGroup string,
) ([]moduleLifecycleTestCase, error) {
	// Creating a SQL server for existing server case only
	serverName := uuid.NewV4().String()
	administratorLogin := generate.NewIdentifier()
	administratorLoginPassword := generate.NewPassword()
	location := "southcentralus"
	createSQLServer := func() error {
		if _, err := armDeployer.Deploy(
			uuid.NewV4().String(),
			resourceGroup,
			location,
			armTemplateNewServerBytes,
			nil, // Go template params
			map[string]interface{}{ // ARM template params
				"serverName":                 serverName,
				"administratorLogin":         administratorLogin,
				"administratorLoginPassword": administratorLoginPassword,
			},
			map[string]string{},
		); err != nil {
			return fmt.Errorf("error deploying ARM template: %s", err)
		}
		return nil
	}

	serverConfig := mssqldb.ServerConfig{
		ServerName:                 serverName,
		ResourceGroupName:          resourceGroup,
		Location:                   location,
		AdministratorLogin:         administratorLogin,
		AdministratorLoginPassword: administratorLoginPassword,
	}
	serverConfigs := []mssqldb.ServerConfig{serverConfig}
	serverConfigsBytes, err := json.Marshal(serverConfigs)
	if err != nil {
		return nil, err
	}
	if err = os.Setenv(
		"AZURE_SQL_SERVERS",
		string(serverConfigsBytes),
	); err != nil {
		return nil, err
	}

	msSQLManager, err := ss.NewManager()
	if err != nil {
		return nil, err
	}
	msSQLConfig, err := mssqldb.GetConfig()
	if err != nil {
		return nil, err
	}

	return []moduleLifecycleTestCase{
		{ // new server scenario
			module:      mssqldb.New(armDeployer, msSQLManager, msSQLConfig),
			description: "new server and database",
			serviceID:   "fb9bc99e-0aa9-11e6-8a8a-000d3a002ed5",
			planID:      "3819fdfa-0aaa-11e6-86f4-000d3a002ed5",
			standardProvisioningContext: service.StandardProvisioningContext{
				Location: "southcentralus",
			},
			provisioningParameters: &mssqldb.ProvisioningParameters{},
			bindingParameters:      &mssqldb.BindingParameters{},
			testCredentials:        testMsSQLCreds(),
		},
		{ // existing server scenario
			module:      mssqldb.New(armDeployer, msSQLManager, msSQLConfig),
			description: "database on an existing server",
			setup:       createSQLServer,
			serviceID:   "fb9bc99e-0aa9-11e6-8a8a-000d3a002ed5",
			planID:      "3819fdfa-0aaa-11e6-86f4-000d3a002ed5",
			standardProvisioningContext: service.StandardProvisioningContext{
				Location: "southcentralus", // This is actually irrelevant for this test
			},
			provisioningParameters: &mssqldb.ProvisioningParameters{
				ServerName: serverName,
			},
			bindingParameters: &mssqldb.BindingParameters{},
			testCredentials:   testMsSQLCreds(),
		},
	}, nil
}

func testMsSQLCreds() func(credentials service.Credentials) error {
	return func(credentials service.Credentials) error {
		cdts, ok := credentials.(*mssqldb.Credentials)
		if !ok {
			return fmt.Errorf("error casting credentials as *mssql.Credentials")
		}

		query := url.Values{}
		query.Add("database", cdts.Database)
		query.Add("encrypt", "true")
		query.Add("TrustServerCertificate", "true")

		u := &url.URL{
			Scheme: "sqlserver",
			User: url.UserPassword(
				cdts.Username,
				cdts.Password,
			),
			Host:     fmt.Sprintf("%s:%d", cdts.Host, cdts.Port),
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
