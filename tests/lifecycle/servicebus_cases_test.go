// +build !unit

package lifecycle

var servicebusTestCases = []serviceLifecycleTestCase{
	{
		group:     "servicebus",
		name:      "servicebus",
		serviceID: "6dc44338-2f13-4bc5-9247-5b1b3c5462d3",
		planID:    "d06817b1-87ea-4320-8942-14b1d060206a",
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"tags": map[string]interface{}{
				"latest-operation": "provision",
			},
		},
	},
}
