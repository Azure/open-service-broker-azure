# [Azure Storage](https://azure.microsoft.com/en-us/services/storage/)

_Note: This module is EXPERIMENTAL. To enable this module, you must run Open Service Broker for Azure with modules.minStability set to `experimental`_

## Services & Plans

### Service: azure-storage

| Plan Name | Description |
|-----------|-------------|
| `general-purpose-v2-storage-account` | This is the second generation storage account and is recommended for most scenarios. This plan provisions a general purpose account only. Create your own containers, files, and tables within this account. |
| `general-purpose-storage-account` | Provisions a general purpose account only. Create your own containers, files, and tables within this account. |
| `blob-storage-account` | Provisions a blob storage account only. Create your own blob containers (only) within this account. |
| `blob-container` | Provisions a blog storage account and a blob container within. |

#### Behaviors

##### Provision

Provisions the storage resources indicated by the applicable plan-- an account
only, or an account with a container.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Y |  |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | Y |  |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |
| ` enableNonHttpsTraffic ` | `string` | Specify whether non-https traffic is enabled. Allowed values:["enabled", "disabled"]. | N | If not provided, "disabled" will be used as the default value. That is, only https traffic is allowed. |
| `accessTier` | `string` | The access tier used for billing.    Allowed values: ["Hot", "Cool"]. Hot storage is optimized for storing data that is accessed frequently ,and cool storage is optimized for storing data that is infrequently accessed and stored for at least 30 days. **Note1** : this field doesn't exist for plan `general-purpose-storage-account`. **Note2** : this field can only set to "Hot" if you use "Premium_LRS" `accountType` | N | If not provided, "Hot" will be used as the default value. |
| `accountType` | `string` | A combination of account kind and   replication strategy. Allowed values: for plan `blob-storage-account` and `blob-container`: [ "Standard_LRS", "Standard_GRS", "Standard_RAGRS"]; for plan `general-purpose-storage-account`: ["Standard_LRS", "Standard_GRS", "Standard_RAGRS", "Premium_LRS"]; for plan `general-purpose-v2-storage-account`: ["Standard_LRS", "Standard_GRS", "Standard_RAGRS", "Standard_ZRS", "Premium_LRS"]. Check [here](https://docs.microsoft.com/en-us/azure/storage/common/storage-redundancy#choosing-a-replication-option) for detailed explanation of replication strategy. | N | If not provided, "Standard_LRS" will be used as the default value for all plans. |

##### Bind

Returns a copy of one shared set of credentials.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and shared credentials:

| Field Name | Type | Description |
|------------|------|-------------|
| `storageAccountName` | `string` | The storage account name. |
| `accessKey` | `string` | A key (password) for accessing the storage account. |
| ` primaryBlobServiceEndPoint ` | `string` | Primary blob service end point. |
| ` primaryTableServiceEndPoint ` | `string` | Primary table service end point. |
| ` primaryFileServiceEndPoint ` | `string` | Primary file service end point. This field only appears in `general-purpose-v2-storage-account` and `general-purpose-storage-account`. |
| ` primaryQueueServiceEndPoint ` | `string` | Primary queue service end point. This field only appears in `general-purpose-v2-storage-account` and `general-purpose-storage-account`. |
| `containerName` | `string` | The name of the container within the storage account. This field only appears in `blob-container`. |

##### Unbind

Does nothing.

##### Update

Updates an existing storage account.

###### Updating parameters

| Parameter Name            | Type                | Description                                                  | Required |
| ------------------------- | ------------------- | ------------------------------------------------------------ | -------- |
| ` enableNonHttpsTraffic ` | `string`            | Specify whether non-https traffic is enabled. Allowed values:["enabled", "disabled"]. | N        |
| `tags`                    | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N        |
| `accessTier`              | `string`            | The access tier used for billing.    Allowed values: ["Hot", "Cool"]. Hot storage is optimized for storing data that is accessed frequently ,and cool storage is optimized for storing data that is infrequently accessed and stored for at least 30 days. **Note1** : this field doesn't exist for plan `general-purpose-storage-account`. **Note2** : this field can only set to "Hot" if you use "Premium_LRS" `accountType` | N        |
| `accountType`             | `string`            | A combination of account kind and   replication strategy. You can only update ["Standard_LRS", "Standard_GRS", "Standard_RAGRS"] accounts to one of ["Standard_LRS", "Standard_GRS", "Standard_RAGRS"]. For "Standard_ZRS" and "Premium_LRS" accounts, they are not updatable. | N        |

##### Deprovision

Deletes the storage account.
