// +build !unit
// +build experimental

package e2e

var aciTestCases = []e2eTestCase{
	{
		group:     "aci",
		name:      "aci",
		serviceID: "451d5d19-4575-4d4a-9474-116f705ecc95",
		planID:    "d48798e2-21db-405b-abc7-aa6f0ff08f6c",
		provisioningParameters: map[string]interface{}{
			"location":      "eastus",
			"resourceGroup": "placeholder",
			"image":         "nginx",
			"memoryInGb":    1.5,
			"cpuCores":      1,
			"ports":         []int{80, 443},
		},
		bind: true,
	},
}
