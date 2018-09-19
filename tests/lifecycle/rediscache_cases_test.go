// +build !unit

package lifecycle

var rediscacheTestCases = []serviceLifecycleTestCase{
	{
		group:     "rediscache",
		name:      "rediscache-basic-provision-and-update",
		serviceID: "0346088a-d4b2-4478-aa32-f18e295ec1d9",
		planID:    "362b3d1b-5b57-4289-80ad-4a15a760c29c",
		provisioningParameters: map[string]interface{}{
			"location":         "southcentralus",
			"skuCapacity":      1,
			"enableNonSslPort": "disabled",
		},
		updatingParameters: map[string]interface{}{
			"skuCapacity":      2,
			"enableNonSslPort": "enabled",
		},
	},
	{
		group:     "rediscache",
		name:      "rediscache-premium-shardCount-provision-and-update",
		serviceID: "0346088a-d4b2-4478-aa32-f18e295ec1d9",
		planID:    "b1057a8f-9a01-423a-bc35-e168d5c04cf0",
		provisioningParameters: map[string]interface{}{
			"location":   "eastus",
			"shardCount": 2,
		},
		updatingParameters: map[string]interface{}{
			"shardCount": 1,
		},
	},
}
