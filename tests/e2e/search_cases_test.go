// +build !unit

package e2e

var searchTestCases = []e2eTestCase{
	{
		group:     "search",
		name:      "search",
		serviceID: "c54902aa-3027-4c5c-8e96-5b3d3b452f7f",
		planID:    "35bd6e80-5ff5-487e-be0e-338aee9321e4",
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
		},
		bind: true,
	},
}
