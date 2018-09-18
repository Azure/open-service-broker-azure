# [Azure Redis Cache](https://azure.microsoft.com/en-us/services/cache/)

_Note: This module is EXPERIMENTAL. To enable this module, you must run Open Service Broker for Azure with modules.minStability set to `experimental`_

## Services & Plans

### Service: azure-rediscache

| Plan Name  | Description                      |
| ---------- | -------------------------------- |
| `basic`    | Basic Tier, default 250MB Cache  |
| `standard` | Standard Tier, default 1GB Cache |
| `premium`  | Premium Tier, default 6GB Cache  |

#### Behaviors

##### Provision

Provisions a new Redis cache.

###### Provisioning Parameters

| Parameter Name      | Type                | Description                                                  | Required                                                     | Default Value                                                |
| ------------------- | ------------------- | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ |
| `location`          | `string`            | The Azure region in which to provision applicable resources. | Required _unless_ an administrator has configured the broker itself with a default location. | The broker's default location, if configured.                |
| `resourceGroup`     | `string`            | The (new or existing) resource group with which to associate new resources. | N                                                            | If an administrator has configured the broker itself with a default resource group and nonde is specified, that default will be applied, otherwise, a new resource group will be created with a UUID as its name. |
| `skuCapacity`       | `integer`           | The size of the Redis cache to deploy.  Valid values: for C (Basic/Standard) family (0, 1, 2, 3, 4, 5, 6), for P (Premium) family (1, 2, 3, 4). They denotes real size  (250MB, 1GB, 2.5 GB, 6 GB, 13 GB, 26 GB, 53GB) and (6 GB, 13 GB, 26 GB, 53GB) respectively. | N                                                            | If not provided, `0` is used for C (Basic/Standard) family; `1` is used for P (Premium) family. |
| `enableNonSslPort ` | `string`            | Specifies whether the non-SSL Redis server port (6379) is enabled. Valid values: (`enabled`, `disabled`) | N                                                            | If not provided, `disabled` is used. That is, you can't use non-SSL Redis server port by default. |
| `tags`              | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N                                                            | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |

##### Bind

Returns a copy of one shared set of credentials.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and shared credentials:

| Field Name | Type     | Description                                       |
| ---------- | -------- | ------------------------------------------------- |
| `host`     | `string` | The fully-qualified address of the Redis cache.   |
| `port`     | `int`    | The port number to connect to on the Redis cache. |
| `password` | `string` | The password for the Redis cache.                 |
| `uri`      | `string` | The connection string for the Redis cache.        |

**Note**: if `enableNonSslPort` is set to `enabled`, then `port` will be `6379` and the scheme will be `redis` in `uri`; if `enableNonSslPort`  is set to `disabled`, then `port` will be `6380` and the scheme will be `rediss` in `uri`, and you can only use rediss to connect to the Redis cache.

##### Unbind

Does nothing.

##### Update

Updates existing Redis cache.

###### Updating parameters

| Parameter Name     | Type      | Description                                                  | Required |
| ------------------ | --------- | ------------------------------------------------------------ | -------- |
| `skuCapacity`      | `integer` | The size of the Redis cache to deploy.  Valid values: for C (Basic/Standard) family (0, 1, 2, 3, 4, 5, 6), for P (Premium) family (1, 2, 3, 4). They denotes real size  (250MB, 1GB, 2.5 GB, 6 GB, 13 GB, 26 GB, 53GB) and (6 GB, 13 GB, 26 GB, 53GB) respectively. **Note**: you can only update from a smaller capacity to a larger capacity, the reverse is not allowed. | N        |
| `enableNonSslPort` | `string`  | Specifies whether the non-ssl Redis server port (6379) is enabled. Valid values: (`enabled`, `disabled`) | N        |

##### Deprovision

Deletes the Redis cache.