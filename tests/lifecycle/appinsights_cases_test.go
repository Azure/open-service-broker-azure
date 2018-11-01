// +build !unit

package lifecycle

var appinsightsTestCases = []serviceLifecycleTestCase{
	{
		group:     "appinsights",
		name:      "appinsights",
		serviceID: "66130ee7-451b-4c61-8b78-d5c426a06f3e",
		planID:    "a75e8854-591a-4ef2-b3f1-b311d2a02902",
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"tags": map[string]interface{}{
				"test": "test",
			},
		},
	},
}
