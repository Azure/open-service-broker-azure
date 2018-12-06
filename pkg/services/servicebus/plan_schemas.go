package servicebus

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/ptr"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func generateNamespaceProvisioningParamsSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		RequiredProperties: []string{"location", "resourceGroup"},
		PropertySchemas: map[string]service.PropertySchema{
			"location": &service.StringPropertySchema{
				Title: "Location",
				Description: "The Azure region in which to provision" +
					" applicable resources.",
				CustomPropertyValidator: azure.LocationValidator,
			},
			"resourceGroup": &service.StringPropertySchema{
				Title: "Resource group",
				Description: "The (new or existing) resource group with which" +
					" to associate new resources.",
			},
			"tags": &service.ObjectPropertySchema{
				Title: "Tags",
				Description: "Tags to be applied to new resources," +
					" specified as key/value pairs.",
				Additional: &service.StringPropertySchema{},
			},
		},
	}
}

func generateQueueProvisioningParamsSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{
			"queueName": &service.StringPropertySchema{
				Title:       "Queue Name",
				Description: "The name of the queue.",
			},
			// TODO: add custom validators to configurations
			"maxQueueSize": &service.IntPropertySchema{
				Title: "Max Queue Size",
				Description: "The maximum size of the queue in megabytes, " +
					"which is the size of memory allocated for the queue.",
				DefaultValue: ptr.ToInt64(1024),
			},
			"messageTimeToLive": &service.StringPropertySchema{
				Title:        "Message Time To Live",
				Description:  "Default message timespan to live value.",
				DefaultValue: "PT336H",
			},
			"lockDuration": &service.StringPropertySchema{
				Title: "Lock Duration",
				Description: "The amount of time that the message is locked " +
					"for other receivers.",
				DefaultValue: "PT30S",
			},
		},
	}
}

func generateTopicProvisioningParamsSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{
			"topicName": &service.StringPropertySchema{
				Title:       "Topic Name",
				Description: "The name of the topic.",
			},
			// TODO: add custom validators to configurations
			"maxTopicSize": &service.IntPropertySchema{
				Title: "Max Topic Size",
				Description: "The maximum size of the topic in megabytes, " +
					"which is the size of memory allocated for the topic.",
				DefaultValue: ptr.ToInt64(1024),
			},
			"messageTimeToLive": &service.StringPropertySchema{
				Title:        "Message Time To Live",
				Description:  "Default message timespan to live value.",
				DefaultValue: "PT336H",
			},
		},
	}
}
