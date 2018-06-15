// +build !unit
// +build experimental

package e2e

import "open-service-broker-azure/pkg/service"

var cosmosdbTestCases = []e2eTestCase{
	{ // SQL API all-in-one scenario
		group:     "cosmosdb",
		name:      "sql-api-all-in-one",
		serviceID: "58d9fbbd-7041-4dbe-aabe-6268cd31de84",
		planID:    "58d7223d-934e-4fb5-a046-0c67781eb24e",
		provisioningParameters: service.CombinedProvisioningParameters{
			"location":      "eastus",
			"resourceGroup": "placeholder",
			"ipFilters": map[string]interface{}{
				"allowedIPRanges": []string{"0.0.0.0/0"},
			},
		},
		bind: true,
	},
	{ // SQL API account only scenario
		group:     "cosmosdb",
		name:      "sql-api-account-only",
		serviceID: "6330de6f-a561-43ea-a15e-b99f44d183e6",
		planID:    "71168d1a-c704-49ff-8c79-214dd3d6f8eb",
		provisioningParameters: map[string]interface{}{
			"location":      "eastus",
			"resourceGroup": "placeholder",
			"alias":         "cosmos-account",
			"ipFilters": map[string]interface{}{
				"allowedIPRanges": []string{"0.0.0.0/0"},
			},
		},
		bind: true,
		childTestCases: []*e2eTestCase{
			{ // SQL API database only scenario
				group:     "cosmosdb",
				name:      "database-only",
				serviceID: "87c5132a-6d76-40c6-9621-0c7b7542571b",
				planID:    "c821c68c-c8e0-4176-8cf2-f0ca582a07a3",
				provisioningParameters: map[string]interface{}{
					"parentAlias": "cosmos-account",
				},
				bind: true,
			},
		},
	},
	{ // Graph API scenario
		group:     "cosmosdb",
		name:      "graph-api-account-only",
		serviceID: "5f5252a0-6922-4a0c-a755-f9be70d7c79b",
		planID:    "126a2c47-11a3-49b1-833a-21b563de6c04",
		provisioningParameters: service.CombinedProvisioningParameters{
			"location":      "eastus",
			"resourceGroup": "placeholder",
			"ipFilters": map[string]interface{}{
				"allowedIPRanges": []string{"0.0.0.0/0"},
			},
			"consistencyPolicy": map[string]interface{}{
				"defaultConsistencyLevel": "BoundedStaleness",
				"boundedStaleness": map[string]interface{}{
					"maxStalenessPrefix":   10,
					"maxIntervalInSeconds": 500,
				},
			},
		},
		bind: true,
	},
	{ // Table API scenario
		group:     "cosmosdb",
		name:      "table-api-account-only",
		serviceID: "37915cad-5259-470d-a7aa-207ba89ada8c",
		planID:    "c970b1e8-794f-4d7c-9458-d28423c08856",
		provisioningParameters: service.CombinedProvisioningParameters{
			"location":      "southcentralus",
			"resourceGroup": "placeholder",
			"ipFilters": map[string]interface{}{
				"allowedIPRanges": []string{"0.0.0.0/0"},
			},
		},
		bind: true,
	},
	{ // MongoDB API scenario
		group:     "cosmosdb",
		name:      "mongo-api-account-only",
		serviceID: "8797a079-5346-4e84-8018-b7d5ea5c0e3a",
		planID:    "86fdda05-78d7-4026-a443-1325928e7b02",
		provisioningParameters: service.CombinedProvisioningParameters{
			"location":      "eastus",
			"resourceGroup": "placeholder",
			"ipFilters": map[string]interface{}{
				"allowedIPRanges": []string{"0.0.0.0/0"},
			},
		},
		bind: true,
	},
}
