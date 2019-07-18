// +build !unit

package lifecycle

import (
	"context"
	"database/sql"
	"fmt"

	networkSDK "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-01-01/network" // nolint: lll
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	_ "github.com/lib/pq" // Postgres SQL driver
	uuid "github.com/satori/go.uuid"
)

var postgresqlDBMSAlias = uuid.NewV4().String()
var postgresqlV10DBMSAlias = uuid.NewV4().String()

var postgresqlTestCases = []serviceLifecycleTestCase{
	{
		group:           "postgresql",
		name:            "all-in-one",
		serviceID:       "b43b4bba-5741-4d98-a10b-17dc5cee0175",
		planID:          "90f27532-0286-42e5-8e23-c3bb37191368",
		testCredentials: testPostgreSQLCreds,
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
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
			"sslEnforcement": "disabled",
			"extensions": []interface{}{
				"uuid-ossp",
				"postgis",
			},
			"backupRedundancy": "geo",
		},
		updatingParameters: map[string]interface{}{
			"cores":           2,
			"storage":         25,
			"backupRetention": 35,
		},
	},
	{
		group:     "postgresql",
		name:      "dbms-only",
		serviceID: "d3f74b44-79bc-4d1e-bf7d-c247c2b851f9",
		planID:    "73191861-04b3-4d0b-a29b-429eb15a83d4",
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"alias":    postgresqlDBMSAlias,
			"firewallRules": []interface{}{
				map[string]interface{}{
					"name":           "AllowAll",
					"startIPAddress": "0.0.0.0",
					"endIPAddress":   "255.255.255.255",
				},
			},
		},
		childTestCases: []*serviceLifecycleTestCase{
			{ // database only scenario
				group:           "postgresql",
				name:            "database-only",
				serviceID:       "25434f16-d762-41c7-bbdd-8045d7f74ca6",
				planID:          "df6f5ef1-e602-406b-ba73-09c107d1e31b",
				testCredentials: testPostgreSQLCreds,
				provisioningParameters: map[string]interface{}{
					"parentAlias": postgresqlDBMSAlias,
					"extensions": []interface{}{
						"uuid-ossp",
						"postgis",
					},
				},
			},
		},
	},
	{
		group:           "postgresql",
		name:            "all-in-one-v10",
		serviceID:       "32d3b4e0-e68f-4e96-93d4-35fd380f0874",
		planID:          "6caf83ec-5cc1-42a0-9b34-0d163d73064c",
		testCredentials: testPostgreSQLCreds,
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
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
			"sslEnforcement": "disabled",
			"extensions": []interface{}{
				"uuid-ossp",
				"postgis",
			},
			"backupRedundancy": "geo",
		},
		preProvisionFns: []preProvisionFn{
			createVirtualNetworkForPostgres,
		},
		updatingParameters: map[string]interface{}{
			"cores":           2,
			"storage":         25,
			"backupRetention": 35,
		},
	},
	{
		group:     "postgresql",
		name:      "dbms-only-v10",
		serviceID: "cabd3125-5a13-46ea-afad-a69582af9578",
		planID:    "f5218659-72ba-4fd3-9567-afd52d871fee",
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"alias":    postgresqlV10DBMSAlias,
			"firewallRules": []interface{}{
				map[string]interface{}{
					"name":           "AllowAll",
					"startIPAddress": "0.0.0.0",
					"endIPAddress":   "255.255.255.255",
				},
			},
		},
		preProvisionFns: []preProvisionFn{
			createVirtualNetworkForPostgres,
		},
		childTestCases: []*serviceLifecycleTestCase{
			{ // database only scenario
				group:           "postgresql",
				name:            "database-only-v10",
				serviceID:       "1fd01042-3b70-4612-ac19-9ced0b2a1525",
				planID:          "672f80d5-8c9e-488f-b9ce-41142c205e99",
				testCredentials: testPostgreSQLCreds,
				provisioningParameters: map[string]interface{}{
					"parentAlias": postgresqlV10DBMSAlias,
					"extensions": []interface{}{
						"uuid-ossp",
						"postgis",
					},
				},
			},
		},
	},
	// Test case for specifying server name, admin username and admin password
	{
		group:           "postgresql",
		name:            "all-in-one-v10-specified-server-info",
		serviceID:       "32d3b4e0-e68f-4e96-93d4-35fd380f0874",
		planID:          "6caf83ec-5cc1-42a0-9b34-0d163d73064c",
		testCredentials: testPostgreSQLCreds,
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
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
			"sslEnforcement": "disabled",
			"extensions": []interface{}{
				"uuid-ossp",
				"postgis",
			},
			"backupRedundancy": "geo",
			"adminAccountSettings": map[string]interface{}{
				"adminUsername": "postgresqladmin",
				"adminPassword": generate.NewPassword(),
			},
			"serverName": uuid.NewV4().String() + "-specified",
		},
		preProvisionFns: []preProvisionFn{
			createVirtualNetworkForPostgres,
		},
		updatingParameters: map[string]interface{}{
			"cores":           2,
			"storage":         25,
			"backupRetention": 35,
		},
	},
	{
		group:     "postgresql",
		name:      "dbms-only-v10-specified-server-info",
		serviceID: "cabd3125-5a13-46ea-afad-a69582af9578",
		planID:    "f5218659-72ba-4fd3-9567-afd52d871fee",
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"alias":    postgresqlV10DBMSAlias + "-2",
			"firewallRules": []interface{}{
				map[string]interface{}{
					"name":           "AllowAll",
					"startIPAddress": "0.0.0.0",
					"endIPAddress":   "255.255.255.255",
				},
			},
			"adminAccountSettings": map[string]interface{}{
				"adminUsername": "postgresqladmin",
				"adminPassword": generate.NewPassword(),
			},
			"serverName": uuid.NewV4().String() + "-specified",
		},
		preProvisionFns: []preProvisionFn{
			createVirtualNetworkForPostgres,
		},
		childTestCases: []*serviceLifecycleTestCase{
			{ // database only scenario
				group:           "postgresql",
				name:            "database-only-v10",
				serviceID:       "1fd01042-3b70-4612-ac19-9ced0b2a1525",
				planID:          "672f80d5-8c9e-488f-b9ce-41142c205e99",
				testCredentials: testPostgreSQLCreds,
				provisioningParameters: map[string]interface{}{
					"parentAlias": postgresqlV10DBMSAlias + "-2",
					"extensions": []interface{}{
						"uuid-ossp",
						"postgis",
					},
				},
			},
		},
	},
}

func testPostgreSQLCreds(credentials map[string]interface{}) error {
	var connectionStrTemplate string
	if credentials["sslRequired"].(bool) {
		connectionStrTemplate =
			"postgres://%s:%s@%s/%s?sslmode=require"
	} else {
		connectionStrTemplate = "postgres://%s:%s@%s/%s"
	}
	db, err := sql.Open("postgres", fmt.Sprintf(
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

// nolint: lll
func createVirtualNetworkForPostgres(
	ctx context.Context,
	resourceGroup string,
	parent *service.Instance,
	pp *map[string]interface{},
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	azureConfig, err := getAzureConfig()
	if err != nil {
		return fmt.Errorf("error getting azure config %s", err)
	}
	authorizer, err := getBearerTokenAuthorizer(azureConfig)
	if err != nil {
		return fmt.Errorf("error getting authorizer %s", err)
	}
	virtualNetworksClient := networkSDK.NewVirtualNetworksClientWithBaseURI(
		azureConfig.Environment.ResourceManagerEndpoint,
		azureConfig.SubscriptionID,
	)
	virtualNetworksClient.Authorizer = authorizer
	virtualNetworkName := uuid.NewV4().String()
	subnetName := "default"
	vnResult, err := virtualNetworksClient.CreateOrUpdate(
		ctx,
		resourceGroup,
		virtualNetworkName,
		networkSDK.VirtualNetwork{
			Location: to.StringPtr("southcentralus"),
			VirtualNetworkPropertiesFormat: &networkSDK.VirtualNetworkPropertiesFormat{
				AddressSpace: &networkSDK.AddressSpace{
					AddressPrefixes: &[]string{"172.19.0.0/16"},
				},
				Subnets: &[]networkSDK.Subnet{
					{
						Name: &subnetName,
						SubnetPropertiesFormat: &networkSDK.SubnetPropertiesFormat{
							AddressPrefix: to.StringPtr("172.19.0.0/24"),
						},
					},
				},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("error creating virtual network %s", err)
	}
	if err = vnResult.WaitForCompletionRef(
		ctx,
		virtualNetworksClient.Client,
	); err != nil {
		return fmt.Errorf("error waiting for virtual network creation complete %s", err)
	}
	subnetID := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s/subnets/%s",
		azureConfig.SubscriptionID,
		resourceGroup,
		virtualNetworkName,
		subnetName,
	)

	(*pp)["virtualNetworkRules"] = []interface{}{
		map[string]interface{}{
			"name":     "test-subnet",
			"subnetId": subnetID,
		},
	}

	return nil
}
