package rediscache

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

type planDetail struct {
	planName          string
	allowedCapacity   []int64
	allowedShardCount []int64
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
		ips.PropertySchemas["shardCount"] = &service.IntPropertySchema{
			Title: "Shard Count",
			Description: "The number of shards to be created on a Premium Cluster Cache. " +
				"This action is irreversible. The number of shards can be changed later.",
			AllowedValues: pd.allowedShardCount,
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
