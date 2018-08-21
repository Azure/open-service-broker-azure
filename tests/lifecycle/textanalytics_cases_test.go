// +build !unit

package lifecycle

var textanalyticsTestCases = []serviceLifecycleTestCase{
	{ // Text analytics free tier
		group:     "textanalytics",
		name:      "text-analytics-free",
		serviceID: "8f6c848a-4ce1-4a69-9248-63545d3e7e9c",
		planID:    "d5a0f91f-10da-42fc-b792-656a616d9ec2",
		provisioningParameters: map[string]interface{}{
			"location": "eastus",
		},
	},
	{ // Text analytics standard-s0 tier
		group:     "textanalytics",
		name:      "text-analytics-s0",
		serviceID: "8f6c848a-4ce1-4a69-9248-63545d3e7e9c",
		planID:    "7f49713b-2689-4c66-bac9-85a024c0fb9e",
		provisioningParameters: map[string]interface{}{
			"location": "eastus",
		},
	},
	{ // Text analytics standard-s1 tier
		group:     "textanalytics",
		name:      "text-analytics-s1",
		serviceID: "8f6c848a-4ce1-4a69-9248-63545d3e7e9c",
		planID:    "55575612-482b-4260-b67e-69be36d83a54",
		provisioningParameters: map[string]interface{}{
			"location": "eastus",
		},
	},
	{ // Text analytics standard-s2 tier
		group:     "textanalytics",
		name:      "text-analytics-s2",
		serviceID: "8f6c848a-4ce1-4a69-9248-63545d3e7e9c",
		planID:    "76bc6a2f-1364-4ef2-8037-d7cfff48f3b6",
		provisioningParameters: map[string]interface{}{
			"location": "eastus",
		},
	},
	{ // Text analytics standard-s3 tier
		group:     "textanalytics",
		name:      "text-analytics-s3",
		serviceID: "8f6c848a-4ce1-4a69-9248-63545d3e7e9c",
		planID:    "1d9a9e7c-80ac-4f23-aabe-876125541f59",
		provisioningParameters: map[string]interface{}{
			"location": "eastus",
		},
	},
	{ // Text analytics standard-s4 tier
		group:     "textanalytics",
		name:      "text-analytics-s4",
		serviceID: "8f6c848a-4ce1-4a69-9248-63545d3e7e9c",
		planID:    "b9db834d-1350-4c50-adaf-f1e59efa2381",
		provisioningParameters: map[string]interface{}{
			"location": "eastus",
		},
	},
}
