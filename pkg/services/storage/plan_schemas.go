package storage

import (
	"regexp"

	"github.com/Azure/open-service-broker-azure/pkg/ptr"
	"github.com/Azure/open-service-broker-azure/pkg/schemas"
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
			"resourceGroup": schemas.GetResourceGroupSchema(),
			"location":      schemas.GetLocationSchema(),
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

	ips.PropertySchemas["accountType"] = &service.StringPropertySchema{
		Title: "Account Type",
		Description: "This field is a combination of account kind and " +
			" replication strategy",
		DefaultValue:  "Standard_LRS",
		AllowedValues: accountTypeMap[serviceName],
	}

	if serviceName != serviceGeneralPurposeV1 {
		ips.PropertySchemas["accessTier"] = &service.StringPropertySchema{
			Title:         "Access Tier",
			Description:   "The access tier used for billing.",
			DefaultValue:  hot,
			AllowedValues: []string{hot, cool},
		}
	}

	if serviceName == serviceBlobAllInOne {
		ips.PropertySchemas["containerName"] = &service.StringPropertySchema{
			Title: "Container Name",
			Description: "The name of the container which will be created inside" +
				"the blob stroage account",
			AllowedPattern: regexp.MustCompile("^[a-z0-9]+(?:-[a-z0-9]+)*$"),
			MinLength:      ptr.ToInt(3),
			MaxLength:      ptr.ToInt(63),
		}
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

	ips.PropertySchemas["accountType"] = &service.StringPropertySchema{
		Title: "Account Type",
		Description: "This field is a combination of account kind and " +
			" replication strategy",
		DefaultValue:  "Standard_LRS",
		AllowedValues: accountTypeMap["update"],
	}

	if serviceName != serviceGeneralPurposeV1 {
		ips.PropertySchemas["accessTier"] = &service.StringPropertySchema{
			Title:         "Access Tier",
			Description:   "The access tier used for billing.",
			AllowedValues: []string{hot, cool},
		}
	}

	return ips
}

// nolint: lll
func generateBlobContainerProvisioningParamsSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{
			"containerName": &service.StringPropertySchema{
				Title: "Container Name",
				Description: "The name of the container which will be created inside" +
					"the blob stroage account",
				AllowedPattern: regexp.MustCompile("^[a-z0-9]+(?:-[a-z0-9]+)*$"),
				MinLength:      ptr.ToInt(3),
				MaxLength:      ptr.ToInt(63),
			},
		},
	}
}
