# [Azure Storage](https://azure.microsoft.com/en-us/services/storage/)

_Note: This module is EXPERIMENTAL and future releases may break the API._

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
| `location` | `string` | The Azure region in which to provision applicable resources. | Required _unless_ an administrator has configured the broker itself with a default location. | The broker's default location, if configured. |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | N | If an administrator has configured the broker itself with a default resource group and nonde is specified, that default will be applied, otherwise, a new resource group will be created with a UUID as its name. |
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
