// +build !unit

package lifecycle

var storageTestCases = []serviceLifecycleTestCase{
	{ // General Purpose V2 Storage Account
		group:     "storage",
		name:      "general-purpose-v2-account",
		serviceID: "2e2fc314-37b6-4587-8127-8f9ee8b33fea",
		planID:    "e19fb0be-dd1f-4ef0-b44f-88832dca1a66",
		provisioningParameters: map[string]interface{}{
			"location":              "westus",
			"enableNonHttpsTraffic": "enabled",
			"tags": map[string]interface{}{
				"latest-operation": "provision",
			},
			"accessTier":  "Hot",
			"accountType": "Standard_ZRS",
		},
		updatingParameters: map[string]interface{}{
			"enableNonHttpsTraffic": "disabled",
			"tags": map[string]interface{}{
				"latest-operation": "update",
			},
			"accessTier": "Cool",
		},
	},
	{ // General Purpose Storage Account
		group:     "storage",
		name:      "general-purpose-account",
		serviceID: "2e2fc314-37b6-4587-8127-8f9ee8b33fea",
		planID:    "6ddf6b41-fb60-4b70-af99-8ecc4896b3cf",
		provisioningParameters: map[string]interface{}{
			"location":              "southcentralus",
			"enableNonHttpsTraffic": "disabled",
			"accountType":           "Premium_LRS",
		},
		updatingParameters: map[string]interface{}{
			"enableNonHttpsTraffic": "enabled",
		},
	},
	{ // Blob Storage Account
		group:     "storage",
		name:      "blob-account",
		serviceID: "2e2fc314-37b6-4587-8127-8f9ee8b33fea",
		planID:    "800a17e1-f20a-463d-a290-20516052f647",
		provisioningParameters: map[string]interface{}{
			"location":              "eastus",
			"enableNonHttpsTraffic": "enabled",
			"accessTier":            "Cool",
			"accountType":           "Standard_LRS",
		},
		updatingParameters: map[string]interface{}{
			"accessTier":  "Hot",
			"accountType": "Standard_RAGRS",
		},
	},
	{ // Blob Storage Account + Blob Container
		group:     "storage",
		name:      "blob-account-with-container",
		serviceID: "2e2fc314-37b6-4587-8127-8f9ee8b33fea",
		planID:    "189d3b8f-8307-4b3f-8c74-03d069237f70",
		provisioningParameters: map[string]interface{}{
			"location":              "southcentralus",
			"enableNonHttpsTraffic": "enabled",
			"accountType":           "Standard_GRS",
		},
		updatingParameters: map[string]interface{}{
			"accountType": "Standard_LRS",
		},
	},
}
