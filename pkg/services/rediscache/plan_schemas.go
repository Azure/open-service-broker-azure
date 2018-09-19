package rediscache

import (
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
		ips.PropertySchemas["redisConfiguration"] = &service.ObjectPropertySchema{
			Title:       "Redis Configuration",
			Description: "All Redis Settings.",
			PropertySchemas: map[string]service.PropertySchema{
				"rdb-backup-enabled": &service.StringPropertySchema{
					Title:         "RDB backup enabled",
					Description:   "Specifies whether RDB backup is enabled.",
					AllowedValues: []string{"enabled", "disabled"},
				},
				"rdb-backup-frequency": &service.IntPropertySchema{
					Title:         "RDB backup frequency",
					Description:   "The frequency doing backup",
					AllowedValues: []int64{15, 30, 60, 360, 720, 1440},
				},
				"rdb-storage-connection-string": &service.StringPropertySchema{
					Title:       "RDB storage connection string",
					Description: "The connnection string of the storage account for backup",
				},
			},
			DefaultValue: map[string]interface{}{},
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

	if pd.planName == premium {
		ips.PropertySchemas["redisConfiguration"] = &service.ObjectPropertySchema{
			Title:       "Redis Configuration",
			Description: "All Redis Settings.",
			PropertySchemas: map[string]service.PropertySchema{
				"rdb-backup-enabled": &service.StringPropertySchema{
					Title:         "RDB backup enabled",
					Description:   "Specifies whether RDB backup is enabled.",
					AllowedValues: []string{"enabled", "disabled"},
				},
				"rdb-backup-frequency": &service.IntPropertySchema{
					Title:         "RDB backup frequency",
					Description:   "The frequency doing backup",
					AllowedValues: []int64{15, 30, 60, 360, 720, 1440},
				},
				"rdb-storage-connection-string": &service.StringPropertySchema{
					Title:       "RDB storage connection string",
					Description: "The connnection string of the storage account for backup",
				},
			},
		}
	}

	return ips
}
