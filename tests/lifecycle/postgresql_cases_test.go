// +build !unit

package lifecycle

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Postgres SQL driver
	uuid "github.com/satori/go.uuid"
)

var postgresqlDBMSAlias = uuid.NewV4().String()

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
			"backupRedundancy" : "geo",
		},
		updatingParameters: map[string]interface{}{
			"cores": 2,
			"storage": 25,
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
