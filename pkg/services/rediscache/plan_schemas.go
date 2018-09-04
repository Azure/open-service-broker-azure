package rediscache

import (
	"fmt"
	"net"

	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

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
				Title:        "Enable non-SSL port",
				Description:  "Specifies whether the non-ssl Redis server port (6379) is enabled.",
				DefaultValue: "enabled",
			},
			"skuCapacity": &service.IntPropertySchema{
				Title:         "Sku capacity",
				Description:   "The size of the Redis cache to deploy.",
				AllowedValues: pd.allowedCapacity,
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
		var defaultBackupFrequency int64 = 60
		ips.PropertySchemas["redisConfiguration"] = &service.ObjectPropertySchema{
			Title:       "Redis Configuration",
			Description: "All Redis Settings.",
			PropertySchemas: map[string]service.PropertySchema{
				"rdb-backup-enabled": &service.StringPropertySchema{
					Title:         "RDB backup enabled",
					Description:   "Specifies whether RDB backup is enabled.",
					AllowedValues: []string{"enabled", "disabled"},
					DefaultValue:  "disabled",
				},
				"rdb-backup-frequency": &service.IntPropertySchema{
					Title:         "RDB backup frequency",
					Description:   "The frequency doing backup",
					AllowedValues: []int64{15, 30, 60, 360, 720, 1440},
					DefaultValue:  &defaultBackupFrequency,
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

func (pd planDetail) getUpdatingParamsSchema() service.InputParametersSchema {
	ips := service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{
			"enableNonSslPort": &service.StringPropertySchema{
				Title:       "Enable non-SSL port",
				Description: "Specifies whether the non-ssl Redis server port (6379) is enabled.",
			},
			"skuCapacity": &service.IntPropertySchema{
				Title:         "Sku capacity",
				Description:   "The size of the Redis cache to deploy.",
				AllowedValues: pd.allowedCapacity,
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
