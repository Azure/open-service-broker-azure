# [Azure Redis Cache](https://azure.microsoft.com/en-us/services/cache/)

_Note: This module is EXPERIMENTAL and is not included in the General Availability release of Open Service Broker for Azure. It will be added in a future OSBA release._

## Services & Plans

### Service: azure-rediscache

| Plan Name | Description |
|-----------|-------------|
| `basic` | Basic Tier, default 250MB Cache |
| `standard` | Standard Tier, default 1GB Cache |
| `premium` | Premium Tier, default 6GB Cache |

#### Behaviors

##### Provision

Provisions a new Redis cache.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Required _unless_ an administrator has configured the broker itself with a default location. | The broker's default location, if configured. |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | N | If an administrator has configured the broker itself with a default resource group and nonde is specified, that default will be applied, otherwise, a new resource group will be created with a UUID as its name. |
| `skuCapacity` | `integer` | The size of the Redis cache to deploy.  Valid values: for C (Basic/Standard) family (0, 1, 2, 3, 4, 5, 6), for P (Premium) family (1, 2, 3, 4). | N | If not provided, `0` is used for C (Basic/Standard) family; `1` is used for P (Premium) family. |
| `enableNonSslPort ` | `string` | Specifies whether the non-ssl Redis server port (6379) is enabled. Valid values: (`enabled`, `disabled`) | N | If not provided, `enabled` is used. **Note**:  this behavior is different from Azure portal. In OSBA, non-SSL port is enabled by default. That's because we want to make sure your application can work normally even your application doesn't support Redis with SSL. |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |

For `premium` plan, following provisioning parameter is available:

| Parameter Name                                           | Type      | Description                                                  | Required                                                     | Default Value                                                |
| -------------------------------------------------------- | --------- | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `shardCount`                                             | `integer` | The number of shards to be created on a Premium Cluster Cache. This action is irreversible. The number of shards can be changed later. | N                                                            | If not specified, no additional shard will be created.       |
| ` subnetId `                                             | `string`  | The full resource ID of a subnet in a virtual network to deploy the Redis cache in. The subnet should be in the same region with Redis cache. Example format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/virtualNetworks/{vn}/subnets/{sn} | N                                                            | If not specified, the Redis cache will be deployed as common. |
| `staticIP`                                               | `string`  | Static IP address. Required when deploying a Redis cache inside an existing Azure Virtual Network. Only valid when `subnetId` is provided. | N                                                            | If `staticIP` **is not** specified and `subnetId` **is** specified, one valid IP will be chosen randomly in the subnet. |
| `redisConfiguration`                                     | `object`  | Redis Settings. See below possible keys.                     | N                                                            | null object                                                  |
| `redisConfiguration`.` rdb-backup-enabled `              | `string`  | Specifies whether RDB backup is enabled. Valid values: (`enabled`, `disabled`) | N                                                            | If not specified, RDB backup will be disabled by default.    |
| `redisConfiguration`.` rdb-backup-frequency `            | `integer` | The frequency doing backup in minutes. Valid values: ( 15, 30, 60, 360, 720, 1440 ) | Yes when ` rdb-backup-enabled ` is set to `enabled`; otherwise is not required. |                                                              |
| `redisConfiguration`.`  rdb-storage-connection-string  ` | `string`  | The connnection string of the storage account for backup.    | Yes when ` rdb-backup-enabled ` is set to `enabled`; otherwise is not required. |                                                              |

##### Bind

Returns a copy of one shared set of credentials.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and shared credentials:

| Field Name | Type | Description |
|------------|------|-------------|
| `host` | `string` | The fully-qualified address of the Redis cache. |
| `port` | `int` | The port number to connect to on the Redis cache. |
| `password` | `string` | The password for the Redis cache. |
| `uri` | `string` | The connection string for the Redis cache. |

**Note**: if `enableNonSslPort` is set to `enabled`, then `port` will be `6379` and the scheme will be `redis` in `uri`; if `enableNonSslPort`  is set to `disabled`, then `port` will be `6380` and the scheme will be `rediss` in `uri`, and you can only use rediss to connect to the Redis cache.

##### Unbind

Does nothing.

##### Update

Updates existing Redis cache.

###### Updating parameters

| Parameter Name     | Type      | Description                                                  | Required |
| ------------------ | --------- | ------------------------------------------------------------ | -------- |
| `skuCapacity`      | `integer` | The size of the Redis cache to deploy.  Valid values: for C (Basic/Standard) family (0, 1, 2, 3, 4, 5, 6), for P (Premium) family (1, 2, 3, 4).  **Note**: you can only update from a smaller capacity to a larger capacity, the reverse is not allowed. | N        |
| `enableNonSslPort` | `string`  | Specifies whether the non-ssl Redis server port (6379) is enabled. Valid values: (`enabled`, `disabled`) | N        |

For `premium` plan, following updating parameter is available:

| Parameter Name                                       | Type      | Description                                                  | Required                                                     |
| ---------------------------------------------------- | --------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `shardCount`                                         | `integer` | The number of shards to be created on a Premium Cluster Cache. This action is irreversible. The number of shards can be changed later. | N                                                            |
| `redisConfiguration`                                 | `object`  | Redis Settings. See below possible keys.                     | N                                                            |
| `redisConfiguration`.`rdb-backup-enabled`            | `string`  | Specifies whether RDB backup is enabled. Valid values: (`enabled`, `disabled`) | N                                                            |
| `redisConfiguration`.`rdb-backup-frequency`          | `integer` | The frequency doing backup in minutes. Valid values: ( 15, 30, 60, 360, 720, 1440 ) | Yes when `rdb-backup-enabled` is set to `enabled`; otherwise is not required. |
| `redisConfiguration`.`rdb-storage-connection-string` | `string`  | The connnection string of the storage account for backup.    | Yes when `rdb-backup-enabled` is set to `enabled`; otherwise is not required. |

##### Deprovision

Deletes the Redis cache.
