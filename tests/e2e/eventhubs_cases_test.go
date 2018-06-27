// +build !unit

package e2e

var eventhubsTestCases = []e2eTestCase{
	{
		group:     "eventhubs",
		name:      "eventhubs",
		serviceID: "7bade660-32f1-4fd7-b9e6-d416d975170b",
		planID:    "80756db5-a20c-495d-ae70-62cf7d196a3c",
		provisioningParameters: map[string]interface{}{
			"location":      "southcentralus",
			"resourceGroup": "placeholder",
		},
		bind: true,
	},
}
