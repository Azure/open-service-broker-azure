// +build !unit

package lifecycle

import (
	"database/sql"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	_ "github.com/go-sql-driver/mysql" // MySQL SQL driver
	uuid "github.com/satori/go.uuid"
)

var mysqlDBMSAlias = uuid.NewV4().String()

var mysqlTestCases = []serviceLifecycleTestCase{
	{
		group:     "mysql",
		name:      "all-in-one",
		serviceID: "997b8372-8dac-40ac-ae65-758b4a5075a5",
		planID:    "eae202c3-521c-46d1-a047-872dacf781fd",
		provisioningParameters: map[string]interface{}{
			"location":       "southcentralus",
			"sslEnforcement": "disabled",
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
			"backupRedundancy": "geo",
		},
		updatingParameters: map[string]interface{}{
			"cores":           2,
			"storage":         25,
			"backupRetention": 35,
		},
		testCredentials: testMySQLCreds,
	},
	{
		group:     "mysql",
		name:      "dbms-only",
		serviceID: "30e7b836-199d-4335-b83d-adc7d23a95c2",
		planID:    "b242a78f-9946-406a-af67-813c56341960",
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"alias":    mysqlDBMSAlias,
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
				group:           "mysql",
				name:            "database-only",
				serviceID:       "6704ae59-3eae-49e9-82b4-4cbcc00edf08",
				planID:          "ec77bd04-2107-408e-8fde-8100c1ce1f46",
				testCredentials: testMySQLCreds,
				provisioningParameters: map[string]interface{}{
					"parentAlias": mysqlDBMSAlias,
				},
			},
		},
	},
	// Test case for specifying server name, admin username and admin password,
	{
		group:     "mysql",
		name:      "all-in-one-specified-info",
		serviceID: "997b8372-8dac-40ac-ae65-758b4a5075a5",
		planID:    "eae202c3-521c-46d1-a047-872dacf781fd",
		provisioningParameters: map[string]interface{}{
			"location":       "southcentralus",
			"sslEnforcement": "disabled",
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
			"backupRedundancy": "geo",
			"adminAccountSettings": map[string]interface{}{
				"adminUsername": "osbaciadmin",
				"adminPassword": generate.NewPassword(),
			},
			"serverName": "osbaciservername",
		},
		updatingParameters: map[string]interface{}{
			"cores":           2,
			"storage":         25,
			"backupRetention": 35,
		},
		testCredentials: testMySQLCreds,
	},
	{
		group:     "mysql",
		name:      "dbms-only-specified-info",
		serviceID: "30e7b836-199d-4335-b83d-adc7d23a95c2",
		planID:    "b242a78f-9946-406a-af67-813c56341960",
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"alias":    mysqlDBMSAlias + "-2",
			"firewallRules": []interface{}{
				map[string]interface{}{
					"name":           "AllowAll",
					"startIPAddress": "0.0.0.0",
					"endIPAddress":   "255.255.255.255",
				},
			},
			"adminAccountSettings": map[string]interface{}{
				"adminUsername": "osbaciadmin",
				"adminPassword": generate.NewPassword(),
			},
			"serverName": "osbaciservername",
		},
		childTestCases: []*serviceLifecycleTestCase{
			{ // database only scenario
				group:           "mysql",
				name:            "database-only",
				serviceID:       "6704ae59-3eae-49e9-82b4-4cbcc00edf08",
				planID:          "ec77bd04-2107-408e-8fde-8100c1ce1f46",
				testCredentials: testMySQLCreds,
				provisioningParameters: map[string]interface{}{
					"parentAlias": mysqlDBMSAlias + "-2",
				},
			},
		},
	},
}

func testMySQLCreds(credentials map[string]interface{}) error {

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
