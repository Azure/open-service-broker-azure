# [Azure Storage](https://azure.microsoft.com/en-us/services/storage/)

_Note: This module is EXPERIMENTAL. To enable this module, you must run Open Service Broker for Azure with modules.minStability set to `experimental`_

## Services & Plans

### Service: azure-storage

| Plan Name | Description |
|-----------|-------------|
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
| `containerName` | `string` | If applicable, the name of the container within the storage account. |

##### Unbind

Does nothing.

##### Deprovision

Deletes the storage account.
