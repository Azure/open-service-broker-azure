package rediscache

import (
	"fmt"
	"net"

	"github.com/Azure/open-service-broker-azure/pkg/schemas"
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
			"resourceGroup": schemas.GetResourceGroupSchema(),
			"location":      schemas.GetLocationSchema(),
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
			"tags": &service.ObjectPropertySchema{
				Title: "Tags",
				Description: "Tags to be applied to new resources," +
					" specified as key/value pairs.",
				Additional: &service.StringPropertySchema{},
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
		ips.PropertySchemas["shardCount"] = &service.IntPropertySchema{
			Title: "Shard Count",
			Description: "The number of shards to be created on a Premium Cluster Cache. " +
				"This action is irreversible. The number of shards can be changed later.",
			AllowedValues: pd.allowedShardCount,
		}
		ips.PropertySchemas["subnetSettings"] = &service.ObjectPropertySchema{
			Title: "Subnet Settings",
			Description: "Setting to deploy the Redis cache inside a subnet, so that the " +
				"cache is only accessible in the subnet",
			DefaultValue: map[string]interface{}{},
			PropertySchemas: map[string]service.PropertySchema{
				"subnetId": &service.StringPropertySchema{
					Title: "Subnet ID",
					Description: "The full resource ID of a subnet in a virtual network to deploy " +
						"the Redis cache in",
					DefaultValue: "",
				},
				"subnetIP": &service.StringPropertySchema{
					Title: "Subnet IP",
					Description: "Static IP address. Required when deploying a Redis cache inside " +
						"an existing Azure Virtual Network.",
					DefaultValue:            "",
					CustomPropertyValidator: ipValidator,
				},
			},
			CustomPropertyValidator: subnetSettingsValidator,
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
			"tags": &service.ObjectPropertySchema{
				Title: "Tags",
				Description: "Tags to be applied to new resources," +
					" specified as key/value pairs.",
				Additional: &service.StringPropertySchema{},
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

func subnetSettingsValidator(
	context string,
	value map[string]interface{},
) error {
	_, idOccured := value["subnetId"]
	_, ipOccured := value["subnetIP"]
	if !idOccured && ipOccured {
		return service.NewValidationError(
			context,
			"subnetIP can be provided only when subnetId is provided",
		)
	}
	return nil
}
