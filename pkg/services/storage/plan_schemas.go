package storage

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

const enabled = "enabled"
const disabled = "disabled"
const hot = "Hot"
const cool = "Cool"

type planDetail struct {
	planName string
}

func (pd planDetail) generateProvisioningParamsSchema() service.InputParametersSchema {
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
			"enableNonHttpsTraffic": &service.StringPropertySchema{
				Title:         "Enable non-https traffic",
				Description:   "Specify whether non-https traffic is enabled",
				DefaultValue:  disabled,
				AllowedValues: []string{enabled, disabled},
			},
			"tags": &service.ObjectPropertySchema{
				Title: "Tags",
				Description: "Tags to be applied to new resources," +
					" specified as key/value pairs.",
				Additional: &service.StringPropertySchema{},
			},
		},
	}
	if pd.planName != generalPurposeV1 {
		ips.PropertySchemas["accessTier"] = &service.StringPropertySchema{
			Title:         "Access Tier",
			Description:   "The access tier used for billing.",
			DefaultValue:  hot,
			AllowedValues: []string{hot, cool},
		}
	}
	return ips
}

func (pd planDetail) generateUpdatingParamsSchema() service.InputParametersSchema {
	ips := service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{
			"enableNonHttpsTraffic": &service.StringPropertySchema{
				Title:         "Enable non-https traffic",
				Description:   "Specify whether non-https traffic is enabled",
				AllowedValues: []string{enabled, disabled},
			},
			"tags": &service.ObjectPropertySchema{
				Title: "Tags",
				Description: "Tags to be applied to new resources," +
					" specified as key/value pairs.",
				Additional: &service.StringPropertySchema{},
			},
		},
	}
	if pd.planName != generalPurposeV1 {
		ips.PropertySchemas["accessTier"] = &service.StringPropertySchema{
			Title:         "Access Tier",
			Description:   "The access tier used for billing.",
			AllowedValues: []string{hot, cool},
		}
	}
	return ips
}
