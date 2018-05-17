// +build !unit

package lifecycle

import "github.com/Azure/open-service-broker-azure/pkg/service"

var aciTestCases = []serviceLifecycleTestCase{
	{
		group:     "aci",
		name:      "aci",
		serviceID: "451d5d19-4575-4d4a-9474-116f705ecc95",
		planID:    "d48798e2-21db-405b-abc7-aa6f0ff08f6c",
		location:  "eastus",
		provisioningParameters: service.CombinedProvisioningParameters{
			"image":      "nginx",
			"memoryInGb": 1.5,
			"cpuCores":   float64(1),
			"ports":      []interface{}{float64(80), float64(443)},
		},
	},
}
