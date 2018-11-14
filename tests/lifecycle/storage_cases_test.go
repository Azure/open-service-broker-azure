// +build !unit

package lifecycle

var storageTestCases = []serviceLifecycleTestCase{
	{ // General Purpose V2 Storage Account
		group:     "storage",
		name:      "general-purpose-v2-account",
		serviceID: "9a3e28fe-8c02-49da-9b35-1b054eb06c95",
		planID:    "bc4f766a-c372-479c-b0b4-bd9d0546b3ef",
		provisioningParameters: map[string]interface{}{
			"location":              "eastus",
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
		name:      "general-purpose-v1-account",
		serviceID: "d10ea062-b627-41e8-a240-543b60030694",
		planID:    "9364d013-3690-4ce5-b0a2-b43d9b970b02",
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
		serviceID: "1a5b4582-29a3-48c5-9cac-511fd8c52756",
		planID:    "98ae02ec-da21-4b09-b5e0-e2f9583d565c",
		provisioningParameters: map[string]interface{}{
			"location":              "eastus",
			"enableNonHttpsTraffic": "enabled",
			"accessTier":            "Cool",
			"accountType":           "Standard_LRS",
			"alias":                 "blobAccount",
		},
		updatingParameters: map[string]interface{}{
			"accessTier":  "Hot",
			"accountType": "Standard_RAGRS",
		},
	},
	{ // Blob Storage Account + Blob Container
		group:     "storage",
		name:      "blob-account-all-in-one",
		serviceID: "d799916e-3faf-4bdf-a48b-bf5012a2d38c",
		planID:    "6c3b587d-0f88-4112-982a-dbe541f30669",
		provisioningParameters: map[string]interface{}{
			"location":              "southcentralus",
			"enableNonHttpsTraffic": "enabled",
			"accountType":           "Standard_GRS",
			"containerName":         "blobContainer",
		},
		updatingParameters: map[string]interface{}{
			"accountType": "Standard_LRS",
		},
	},
}
