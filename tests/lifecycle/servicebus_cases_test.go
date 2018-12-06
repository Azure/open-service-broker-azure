// +build !unit

package lifecycle

var serviceBusNamespaceAlias = "test-servicebus-namespace"
var servicebusTestCases = []serviceLifecycleTestCase{
	{
		group:     "servicebus",
		name:      "servicebus-namespace",
		serviceID: "6dc44338-2f13-4bc5-9247-5b1b3c5462d3",
		planID:    "6be0d8b5-381f-4d68-bdfd-a131425d3835",
		provisioningParameters: map[string]interface{}{
			"location": "southcentralus",
			"alias":    serviceBusNamespaceAlias,
			"tags": map[string]interface{}{
				"latest-operation": "provision",
			},
		},
		childTestCases: []*serviceLifecycleTestCase{
			{ // queue scenario
				group:     "servicebus",
				name:      "servicebus-queue",
				serviceID: "0e93fbb8-7904-43a5-82db-81c7d3886a24",
				planID:    "89440eec-a888-49ce-b392-60c653d7a98b",
				provisioningParameters: map[string]interface{}{
					"parentAlias":       serviceBusNamespaceAlias,
					"queueName":         "testqueue",
					"maxQueueSize":      2048,
					"messageTimeToLive": "PT24H23M22S",
					"lockDuration":      "PT2M30S",
				},
			},
			{ // topic scenario
				group:     "servicebus",
				name:      "servicebus-topic",
				serviceID: "dc6d1545-4391-4c7e-ac7e-a8463787fb93",
				planID:    "dd1e4d44-58be-4f34-84ff-f73ccef405e5",
				provisioningParameters: map[string]interface{}{
					"parentAlias":       serviceBusNamespaceAlias,
					"topicName":         "testtopic",
					"maxTopicSize":      4096,
					"messageTimeToLive": "PT276H13M14S",
				},
			},
		},
	},
}
