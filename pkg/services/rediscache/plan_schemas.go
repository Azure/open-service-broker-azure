package rediscache

import (
	"fmt"
	"net"

	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type planDetail struct {
	planName        string
	allowedCapacity []int64
}

// nolint: lll
func (pd planDetail) getProvisioningParamsSchema() service.InputParametersSchema {
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
			"enableNonSslPort": &service.StringPropertySchema{
				Title:         "Enable non-SSL port",
				Description:   "Specifies whether the non-ssl Redis server port (6379) is enabled.",
				AllowedValues: []string{"enabled", "disabled"},
				DefaultValue:  "disabled",
			},
			"skuCapacity": &service.IntPropertySchema{
				Title:         "SKU capacity",
				Description:   "The size of the Redis cache to deploy.",
				AllowedValues: pd.allowedCapacity,
				DefaultValue:  &(pd.allowedCapacity[0]),
			},
		},
	}

	if pd.planName == premium {
		ips.PropertySchemas["subnetId"] = &service.StringPropertySchema{
			Title: "Subnet ID",
			Description: "The full resource ID of a subnet in a virtual network to deploy " +
				"the Redis cache in",
			DefaultValue: "",
		}
		ips.PropertySchemas["subnetIP"] = &service.StringPropertySchema{
			Title: "Subnet IP",
			Description: "Static IP address. Required when deploying a Redis cache inside " +
				"an existing Azure Virtual Network.",
			DefaultValue:            "",
			CustomPropertyValidator: ipValidator,
		}
	}
	return ips
}

// nolint: lll
func (pd planDetail) getUpdatingParamsSchema() service.InputParametersSchema {
	ips := service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{
			"enableNonSslPort": &service.StringPropertySchema{
				Title:         "Enable non-SSL port",
				Description:   "Specifies whether the non-ssl Redis server port (6379) is enabled.",
				AllowedValues: []string{"enabled", "disabled"},
			},
			"skuCapacity": &service.IntPropertySchema{
				Title:         "SKU capacity",
				Description:   "The size of the Redis cache to deploy.",
				AllowedValues: pd.allowedCapacity,
			},
		},
	}
	return ips
}

func ipValidator(context, value string) error {
	ip := net.ParseIP(value)
	if ip == nil {
		return service.NewValidationError(
			context,
			fmt.Sprintf(`"%s" is not a valid IP address`, value),
		)
	}
	return nil
}
