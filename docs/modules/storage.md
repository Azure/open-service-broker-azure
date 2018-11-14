# [Azure Storage](https://azure.microsoft.com/en-us/services/storage/)

_Note: This module is EXPERIMENTAL. To enable this module, you must run Open Service Broker for Azure with modules.minStability set to `experimental`_

## Services & Plans

### Service: azure-storage-general-purpose-v2-storage-account

| Plan Name | Description                                                  |
| --------- | ------------------------------------------------------------ |
| `account` | This plan provisions a general purpose v2 account. General-purpose v2 storage accounts support the latest Azure Storage features and incorporate all of the functionality of general-purpose v1 and Blob storage accounts. General-purpose v2 accounts deliver the lowest per-gigabyte capacity prices for Azure Storage, as well as industry-competitive transaction prices. |

#### Behaviors

##### Provision

Provisions a general purpose v2 storage account.

###### Provisioning Parameters

| Parameter Name          | Type                | Description                                                  | Required | Default Value                                                |
| ----------------------- | ------------------- | ------------------------------------------------------------ | -------- | ------------------------------------------------------------ |
| `location`              | `string`            | The Azure region in which to provision applicable resources. | Y        |                                                              |
| `resourceGroup`         | `string`            | The (new or existing) resource group with which to associate new resources. | Y        |                                                              |
| `enableNonHttpsTraffic` | `string`            | Specify whether non-https traffic is enabled. Allowed values:["enabled", "disabled"]. | N        | If not provided, "disabled" will be used as the default value. That is, only https traffic is allowed. |
| `accessTier`            | `string`            | The access tier used for billing.    Allowed values: ["Hot", "Cool"]. Hot storage is optimized for storing data that is accessed frequently ,and cool storage is optimized for storing data that is infrequently accessed and stored for at least 30 days. **Note** : `accountType` "Premium_LRS" only supports "Hot" in this field | N        | If not provided, "Hot" will be used as the default value.    |
| `accountType`           | `string`            | A combination of account kind and   replication strategy. All possible values: ["Standard_LRS", "Standard_GRS", "Standard_RAGRS", "Standard_ZRS", "Premium_LRS"]. **Note**: ZRS is only available in several regions, check [here](https://docs.microsoft.com/en-us/azure/storage/common/storage-redundancy-zrs#support-coverage-and-regional-availability) for allowed regions to use ZRS. | N        | If not provided, "Standard_LRS" will be used as the default value for all plans. |
| `tags`                  | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N        | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |

##### Bind

Returns a copy of one shared set of credentials.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and shared credentials:

| Field Name                    | Type     | Description                                         |
| ----------------------------- | -------- | --------------------------------------------------- |
| `storageAccountName`          | `string` | The storage account name.                           |
| `accessKey`                   | `string` | A key (password) for accessing the storage account. |
| `primaryBlobServiceEndPoint`  | `string` | Primary blob service end point.                     |
| `primaryTableServiceEndPoint` | `string` | Primary table service end point.                    |
| `primaryFileServiceEndPoint`  | `string` | Primary file service end point.                     |
| `primaryQueueServiceEndPoint` | `string` | Primary queue service end point.                    |

##### Unbind

Does nothing.

##### Update

Updates an existing storage account.

###### Updating parameters

| Parameter Name            | Type                | Description                                                  | Required |
| ------------------------- | ------------------- | ------------------------------------------------------------ | -------- |
| ` enableNonHttpsTraffic ` | `string`            | Specify whether non-https traffic is enabled. Allowed values:["enabled", "disabled"]. | N        |
| `accessTier`              | `string`            | The access tier used for billing.    Allowed values: ["Hot", "Cool"]. Hot storage is optimized for storing data that is accessed frequently ,and cool storage is optimized for storing data that is infrequently accessed and stored for at least 30 days. **Note** : `accountType` "Premium_LRS" only supports "Hot" in this field. | N        |
| `accountType`             | `string`            | A combination of account kind and   replication strategy. You can only update ["Standard_LRS", "Standard_GRS", "Standard_RAGRS"] accounts to one of ["Standard_LRS", "Standard_GRS", "Standard_RAGRS"]. For "Standard_ZRS" and "Premium_LRS" accounts, they are not updatable. | N        |
| `tags`                    | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N        |

##### Deprovision

Deletes the storage account.



### Service: azure-storage-general-purpose-v1-storage-account

| Plan Name | Description                                                  |
| --------- | ------------------------------------------------------------ |
| `account` | This plan provisions a general purpose v1 account. General-purpose v1 accounts provide access to all Azure Storage services, but may not have the latest features or the lowest per gigabyte pricing. |

#### Behaviors

##### Provision

Provisions a general purpose v1 storage account.

###### Provisioning Parameters

| Parameter Name          | Type                | Description                                                  | Required | Default Value                                                |
| ----------------------- | ------------------- | ------------------------------------------------------------ | -------- | ------------------------------------------------------------ |
| `location`              | `string`            | The Azure region in which to provision applicable resources. | Y        |                                                              |
| `resourceGroup`         | `string`            | The (new or existing) resource group with which to associate new resources. | Y        |                                                              |
| `enableNonHttpsTraffic` | `string`            | Specify whether non-https traffic is enabled. Allowed values:["enabled", "disabled"]. | N        | If not provided, "disabled" will be used as the default value. That is, only https traffic is allowed. |
| `accountType`           | `string`            | A combination of account kind and   replication strategy. All possible values: ["Standard_LRS", "Standard_GRS", "Standard_RAGRS", "Premium_LRS"]. | N        | If not provided, "Standard_LRS" will be used as the default value for all plans. |
| `tags`                  | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N        | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |

##### Bind

Returns a copy of one shared set of credentials.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and shared credentials:

| Field Name                    | Type     | Description                                         |
| ----------------------------- | -------- | --------------------------------------------------- |
| `storageAccountName`          | `string` | The storage account name.                           |
| `accessKey`                   | `string` | A key (password) for accessing the storage account. |
| `primaryBlobServiceEndPoint`  | `string` | Primary blob service end point.                     |
| `primaryTableServiceEndPoint` | `string` | Primary table service end point.                    |
| `primaryFileServiceEndPoint`  | `string` | Primary file service end point.                     |
| `primaryQueueServiceEndPoint` | `string` | Primary queue service end point.                    |

##### Unbind

Does nothing.

##### Update

Updates an existing storage account.

###### Updating parameters

| Parameter Name            | Type                | Description                                                  | Required |
| ------------------------- | ------------------- | ------------------------------------------------------------ | -------- |
| ` enableNonHttpsTraffic ` | `string`            | Specify whether non-https traffic is enabled. Allowed values:["enabled", "disabled"]. | N        |
| `accountType`             | `string`            | A combination of account kind and   replication strategy. You can only update ["Standard_LRS", "Standard_GRS", "Standard_RAGRS"] accounts to one of ["Standard_LRS", "Standard_GRS", "Standard_RAGRS"]. For "Premium_LRS" accounts, they are not updatable. | N        |
| `tags`                    | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N        |

##### Deprovision

Deletes the storage account.



### Service: azure-storage-blob-storage-account-and-container

| Plan Name    | Description                                                  |
| ------------ | ------------------------------------------------------------ |
| `all-in-one` | This plan provisions a a specialized Azure storage account for storing block blobs and append blobs, and automatically provisions a blob container within the account. |

#### Behaviors

##### Provision

Provisions a blob storage account and create a container within the account.

###### Provisioning Parameters

| Parameter Name          | Type                | Description                                                  | Required | Default Value                                                |
| ----------------------- | ------------------- | ------------------------------------------------------------ | -------- | ------------------------------------------------------------ |
| `location`              | `string`            | The Azure region in which to provision applicable resources. | Y        |                                                              |
| `resourceGroup`         | `string`            | The (new or existing) resource group with which to associate new resources. | Y        |                                                              |
| `enableNonHttpsTraffic` | `string`            | Specify whether non-https traffic is enabled. Allowed values:["enabled", "disabled"]. | N        | If not provided, "disabled" will be used as the default value. That is, only https traffic is allowed. |
| `accessTier`            | `string`            | The access tier used for billing.    Allowed values: ["Hot", "Cool"]. Hot storage is optimized for storing data that is accessed frequently ,and cool storage is optimized for storing data that is infrequently accessed and stored for at least 30 days. | N        | If not provided, "Hot" will be used as the default value.    |
| `accountType`           | `string`            | A combination of account kind and   replication strategy. All possible values: ["Standard_LRS", "Standard_GRS", "Standard_RAGRS"]. | N        | If not provided, "Standard_LRS" will be used as the default value for all plans. |
| `tags`                  | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N        | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |

##### Bind

Returns a copy of one shared set of credentials.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and shared credentials:

| Field Name                   | Type     | Description                                           |
| ---------------------------- | -------- | ----------------------------------------------------- |
| `storageAccountName`         | `string` | The storage account name.                             |
| `accessKey`                  | `string` | A key (password) for accessing the storage account.   |
| `primaryBlobServiceEndPoint` | `string` | Primary blob service end point.                       |
| `containerName`              | `string` | The name of the container within the storage account. |

##### Unbind

Does nothing.

##### Update

Updates an existing storage account.

###### Updating parameters

| Parameter Name            | Type                | Description                                                  | Required |
| ------------------------- | ------------------- | ------------------------------------------------------------ | -------- |
| ` enableNonHttpsTraffic ` | `string`            | Specify whether non-https traffic is enabled. Allowed values:["enabled", "disabled"]. | N        |
| `accessTier`              | `string`            | The access tier used for billing.    Allowed values: ["Hot", "Cool"]. Hot storage is optimized for storing data that is accessed frequently ,and cool storage is optimized for storing data that is infrequently accessed and stored for at least 30 days. | N        |
| `accountType`             | `string`            | A combination of account kind and   replication strategy.    | N        |
| `tags`                    | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N        |

##### Deprovision

Deletes the storage account and the blob container inside it.



### Service: azure-storage-blob-storage-account

| Plan Name | Description                                                  |
| --------- | ------------------------------------------------------------ |
| `account` | This plan provisions a a specialized Azure storage account for storing block blobs and append blobs. |

#### Behaviors

##### Provision

Provisions a blob storage account.

###### Provisioning Parameters

| Parameter Name          | Type                | Description                                                  | Required | Default Value                                                |
| ----------------------- | ------------------- | ------------------------------------------------------------ | -------- | ------------------------------------------------------------ |
| `location`              | `string`            | The Azure region in which to provision applicable resources. | Y        |                                                              |
| `resourceGroup`         | `string`            | The (new or existing) resource group with which to associate new resources. | Y        |                                                              |
| `enableNonHttpsTraffic` | `string`            | Specify whether non-https traffic is enabled. Allowed values:["enabled", "disabled"]. | N        | If not provided, "disabled" will be used as the default value. That is, only https traffic is allowed. |
| `accessTier`            | `string`            | The access tier used for billing.    Allowed values: ["Hot", "Cool"]. Hot storage is optimized for storing data that is accessed frequently ,and cool storage is optimized for storing data that is infrequently accessed and stored for at least 30 days. | N        | If not provided, "Hot" will be used as the default value.    |
| `accountType`           | `string`            | A combination of account kind and   replication strategy. All possible values: ["Standard_LRS", "Standard_GRS", "Standard_RAGRS"]. | N        | If not provided, "Standard_LRS" will be used as the default value for all plans. |
| `tags`                  | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N        | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |

##### Bind

Returns a copy of one shared set of credentials.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and shared credentials:

| Field Name                   | Type     | Description                                         |
| ---------------------------- | -------- | --------------------------------------------------- |
| `storageAccountName`         | `string` | The storage account name.                           |
| `accessKey`                  | `string` | A key (password) for accessing the storage account. |
| `primaryBlobServiceEndPoint` | `string` | Primary blob service end point.                     |

##### Unbind

Does nothing.

##### Update

Updates an existing storage account.

###### Updating parameters

| Parameter Name            | Type                | Description                                                  | Required |
| ------------------------- | ------------------- | ------------------------------------------------------------ | -------- |
| ` enableNonHttpsTraffic ` | `string`            | Specify whether non-https traffic is enabled. Allowed values:["enabled", "disabled"]. | N        |
| `accessTier`              | `string`            | The access tier used for billing.    Allowed values: ["Hot", "Cool"]. Hot storage is optimized for storing data that is accessed frequently ,and cool storage is optimized for storing data that is infrequently accessed and stored for at least 30 days. | N        |
| `accountType`             | `string`            | A combination of account kind and   replication strategy.    | N        |
| `tags`                    | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N        |

##### Deprovision

Deletes the storage account and the blob container inside it.