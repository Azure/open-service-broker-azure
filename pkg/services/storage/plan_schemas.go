package storage

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

const (
	enabled  = "enabled"
	disabled = "disabled"
	hot      = "Hot"
	cool     = "Cool"
)

// nolint: lll
var accountTypeMap = map[string][]string{
	"update":                {"Standard_LRS", "Standard_GRS", "Standard_RAGRS"},
	serviceBlobAccount:      {"Standard_LRS", "Standard_GRS", "Standard_RAGRS"},
	serviceBlobAllInOne:     {"Standard_LRS", "Standard_GRS", "Standard_RAGRS"},
	serviceGeneralPurposeV1: {"Standard_LRS", "Standard_GRS", "Standard_RAGRS", "Premium_LRS"},
	serviceGeneralPurposeV2: {"Standard_LRS", "Standard_GRS", "Standard_RAGRS", "Premium_LRS", "Standard_ZRS"},
}

// nolint: lll
func generateProvisioningParamsSchema(serviceName string) service.InputParametersSchema {
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
	if serviceName != serviceGeneralPurposeV1 {
		ips.PropertySchemas["accessTier"] = &service.StringPropertySchema{
			Title:         "Access Tier",
			Description:   "The access tier used for billing.",
			DefaultValue:  hot,
			AllowedValues: []string{hot, cool},
		}
	}

	ips.PropertySchemas["accountType"] = &service.StringPropertySchema{
		Title: "Account Type",
		Description: "This field is a combination of account kind and " +
			" replication strategy",
		DefaultValue:  "Standard_LRS",
		AllowedValues: accountTypeMap[serviceName],
	}

	return ips
}

// nolint: lll
func generateUpdatingParamsSchema(serviceName string) service.InputParametersSchema {
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
	if serviceName != serviceGeneralPurposeV1 {
		ips.PropertySchemas["accessTier"] = &service.StringPropertySchema{
			Title:         "Access Tier",
			Description:   "The access tier used for billing.",
			AllowedValues: []string{hot, cool},
		}
	}

	ips.PropertySchemas["accountType"] = &service.StringPropertySchema{
		Title: "Account Type",
		Description: "This field is a combination of account kind and " +
			" replication strategy",
		DefaultValue:  "Standard_LRS",
		AllowedValues: accountTypeMap["update"],
	}

	return ips
}
