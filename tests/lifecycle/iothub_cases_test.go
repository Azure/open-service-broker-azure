// +build !unit

package lifecycle

var iotHubTestCases = []serviceLifecycleTestCase{
	{ // IoT Hub free tier
		group:     "iothub",
		name:      "iot-hub-free",
		serviceID: "afd72c3b-6c2d-40f2-ad0d-d90467989be5",
		planID:    "4d6c40dd-7525-4260-8e4d-f65818197c2b",
		provisioningParameters: map[string]interface{}{
			"location": "eastus",
			"tags": map[string]interface{}{
				"latest-operation": "provision",
			},
		},
	},
	{ // IoT Hub basic b1 tier
		group:     "iothub",
		name:      "iot-hub-basic-b1",
		serviceID: "afd72c3b-6c2d-40f2-ad0d-d90467989be5",
		planID:    "bdff693c-39cb-4590-b4ce-d1a17fab5848",
		provisioningParameters: map[string]interface{}{
			"location":       "eastus",
			"units":          2,
			"partitionCount": 8,
		},
	},
	{ // IoT Hub standard s1 tier
		group:     "iothub",
		name:      "iot-hub-standard-s1",
		serviceID: "afd72c3b-6c2d-40f2-ad0d-d90467989be5",
		planID:    "0dde7e80-1f32-470d-ba0b-9db4fe1826be",
		provisioningParameters: map[string]interface{}{
			"location":       "eastus",
			"units":          1,
			"partitionCount": 2,
		},
	},
}
