package storage

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

const enabled = "enabled"
const disabled = "disabled"

func generateProvisioningParamsSchema() service.InputParametersSchema {
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
			"enableNonHttpsTraffic": &service.StringPropertySchema{
				Title:         "Enable non-https traffic",
				Description:   "Specify whether non-https traffic is enabled",
				DefaultValue:  "disabled",
				AllowedValues: []string{enabled, disabled},
			},
		},
	}
}

func generateUpdatingParamsSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{
			"enableNonHttpsTraffic": &service.StringPropertySchema{
				Title:         "Enable non-https traffic",
				Description:   "Specify whether non-https traffic is enabled",
				AllowedValues: []string{enabled, disabled},
			},
		},
	}
}
