// +build !unit

package e2e

import uuid "github.com/satori/go.uuid"

func getPostgreSQLTestCases() []e2eTestCase {
	alias := uuid.NewV4().String()
	return []e2eTestCase{
		{
			group:     "postgresql",
			name:      "all-in-one",
			serviceID: "b43b4bba-5741-4d98-a10b-17dc5cee0175",
			planID:    "90f27532-0286-42e5-8e23-c3bb37191368",
			provisioningParameters: map[string]interface{}{
				"location": "southcentralus",
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
				"sslEnforcement": "disabled",
				"extensions": []string{
					"uuid-ossp",
					"postgis",
				},
			},
			bind: true,
		},
		{
			group:     "postgresql",
			name:      "dbms-only",
			serviceID: "d3f74b44-79bc-4d1e-bf7d-c247c2b851f9",
			planID:    "73191861-04b3-4d0b-a29b-429eb15a83d4",
			provisioningParameters: map[string]interface{}{
				"alias":    alias,
				"location": "eastus",
				"firewallRules": []map[string]string{
					{
						"name":           "AllowAll",
						"startIPAddress": "0.0.0.0",
						"endIPAddress":   "255.255.255.255",
					},
				},
			},
			childTestCases: []*e2eTestCase{
				{ // database only scenario
					group:     "postgresql",
					name:      "database-only",
					serviceID: "25434f16-d762-41c7-bbdd-8045d7f74ca6",
					planID:    "df6f5ef1-e602-406b-ba73-09c107d1e31b",
					provisioningParameters: map[string]interface{}{
						"parentAlias": alias,
						"extensions": []string{
							"uuid-ossp",
							"postgis",
						},
					},
					bind: true,
				},
			},
		},
	}
}
