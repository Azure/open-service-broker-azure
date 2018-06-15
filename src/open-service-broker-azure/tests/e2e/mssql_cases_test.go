// +build !unit

package e2e

import uuid "github.com/satori/go.uuid"

func getMSSQLTestCases() []e2eTestCase {
	alias := uuid.NewV4().String()
	return []e2eTestCase{
		{ // all-in-one scenario (dtu-based)
			group:     "mssql",
			name:      "all-in-one",
			serviceID: "fb9bc99e-0aa9-11e6-8a8a-000d3a002ed5",
			planID:    "3819fdfa-0aaa-11e6-86f4-000d3a002ed5",
			provisioningParameters: map[string]interface{}{
				"location":      "southcentralus",
				"resourceGroup": "placeholder",
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
			bind: true,
		},
		{ // dbms only scenario
			group:     "mssql",
			name:      "dbms-only",
			serviceID: "a7454e0e-be2c-46ac-b55f-8c4278117525",
			planID:    "24f0f42e-1ab3-474e-a5ca-b943b2c48eee",
			provisioningParameters: map[string]interface{}{
				"alias":         alias,
				"location":      "southcentralus",
				"resourceGroup": "placeholder",
				"firewallRules": []map[string]string{
					{
						"name":           "AllowAll",
						"startIPAddress": "0.0.0.0",
						"endIPAddress":   "255.255.255.255",
					},
				},
			},
			childTestCases: []*e2eTestCase{
				{ // db only scenario (dtu-based)
					group:     "mssql",
					name:      "database-only (DTU)",
					serviceID: "2bbc160c-e279-4757-a6b6-4c0a4822d0aa",
					planID:    "8fa8d759-c142-45dd-ae38-b93482ddc04a",
					provisioningParameters: map[string]interface{}{
						"parentAlias": alias,
					},
					bind: true,
				},
				{ // db only scenario (vcore-based)
					group:     "mssql",
					name:      "database-only (vCore)",
					serviceID: "2bbc160c-e279-4757-a6b6-4c0a4822d0aa",
					planID:    "da591616-77a1-4df8-a493-6c119649bc6b",
					provisioningParameters: map[string]interface{}{
						"parentAlias": alias,
						"cores":       2,
						"storage":     10,
					},
				},
			},
		},
		{ // all-in-one scenario (vcore-based)
			group:     "mssql",
			name:      "all-in-one (vCore)",
			serviceID: "fb9bc99e-0aa9-11e6-8a8a-000d3a002ed5",
			planID:    "c77e86af-f050-4457-a2ff-2b48451888f3",
			provisioningParameters: map[string]interface{}{
				"location":      "southcentralus",
				"resourceGroup": "placeholder",
				"cores":         4,
				"storage":       25,
				"firewallRules": []interface{}{
					map[string]interface{}{
						"name":           "AllowAll",
						"startIPAddress": "0.0.0.0",
						"endIPAddress":   "255.255.255.255",
					},
				},
			},
		},
	}
}
