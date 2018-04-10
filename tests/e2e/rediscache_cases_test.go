// +build !unit

package e2e

var rediscacheTestCases = []e2eTestCase{
	{
		group:     "rediscache",
		name:      "rediscache",
		serviceID: "0346088a-d4b2-4478-aa32-f18e295ec1d9",
		planID:    "362b3d1b-5b57-4289-80ad-4a15a760c29c",
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
		},
		bind: true,
	},
}
