package iothub

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/ptr"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

// nolint: lll
func generateProvisioningParamsSchema(planName string) service.InputParametersSchema {
	ips := service.InputParametersSchema{
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
	if planName == planF1 {
		return ips
	}

	if planName == planS3 {
		ips.PropertySchemas["units"] = &service.IntPropertySchema{
			Title: "Units",
			Description: "Number of IoT hub units. Each IoT Hub is provisioned " +
				"with a certain number of units in a specific tier. " +
				"The tier and number of units determine the maximum " +
				"daily quota of messages that you can send.",
			DefaultValue: ptr.ToInt64(1),
			MinValue:     ptr.ToInt64(1),
			MaxValue:     ptr.ToInt64(10),
		}
	} else {
		ips.PropertySchemas["units"] = &service.IntPropertySchema{
			Title: "Units",
			Description: "Number of IoT hub units. Each IoT Hub is provisioned " +
				"with a certain number of units in a specific tier. " +
				"The tier and number of units determine the maximum " +
				"daily quota of messages that you can send.",
			DefaultValue: ptr.ToInt64(1),
			MinValue:     ptr.ToInt64(1),
			MaxValue:     ptr.ToInt64(200),
		}
	}

	if planName == planB1 || planName == planB2 || planName == planB3 {
		ips.PropertySchemas["partitionCount"] = &service.IntPropertySchema{
			Title: "partitionCount",
			Description: "The number of partitions relates the device-to-cloud " +
				"messages to the number of simultaneous readers of these messages. " +
				"Most IoT hubs only need four partitions.",
			DefaultValue: ptr.ToInt64(4),
			MinValue:     ptr.ToInt64(2),
			MaxValue:     ptr.ToInt64(8),
		}
	} else {
		ips.PropertySchemas["partitionCount"] = &service.IntPropertySchema{
			Title: "partitionCount",
			Description: "The number of partitions relates the device-to-cloud " +
				"messages to the number of simultaneous readers of these messages. " +
				"Most IoT hubs only need four partitions.",
			DefaultValue: ptr.ToInt64(4),
			MinValue:     ptr.ToInt64(2),
			MaxValue:     ptr.ToInt64(32),
		}
	}

	return ips
}
