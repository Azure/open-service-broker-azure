// +build !unit

package e2e

import uuid "github.com/satori/go.uuid"

func getMySQLTestCases() []e2eTestCase {
	alias := uuid.NewV4().String()
	return []e2eTestCase{
		{
			group:     "mysql",
			name:      "all-in-one",
			serviceID: "997b8372-8dac-40ac-ae65-758b4a5075a5",
			planID:    "eae202c3-521c-46d1-a047-872dacf781fd",
			provisioningParameters: map[string]interface{}{
				"location":       "southcentralus",
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
			bind: true,
		},
		{
			group:     "mysql",
			name:      "dbms-only",
			serviceID: "30e7b836-199d-4335-b83d-adc7d23a95c2",
			planID:    "b242a78f-9946-406a-af67-813c56341960",
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
				{
					group:     "mysql",
					name:      "database-only",
					serviceID: "6704ae59-3eae-49e9-82b4-4cbcc00edf08",
					planID:    "ec77bd04-2107-408e-8fde-8100c1ce1f46",
					provisioningParameters: map[string]interface{}{
						"parentAlias": alias,
					},
					bind: true,
				},
			},
		},
	}
}
