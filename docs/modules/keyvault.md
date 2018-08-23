# [Azure Key Vault](https://azure.microsoft.com/en-us/services/key-vault/)

_Note: This module is EXPERIMENTAL. To enable this module, you must run Open Service Broker for Azure with Minimum Stability set to `experimental`_

## Services & Plans

### Service: azure-keyvault

| Plan Name | Description |
|-----------|-------------|
| `standard` | Standard Tier |
| `premium` | Premium Tier |

#### Behaviors

##### Provision
  
Provisions a new Key Vault. The new vault will be named using a new UUID.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `clientId` | `string` | Client ID (username) for an existing service principal, which will be granted access to the new vault.| Y | |
| `clientSecret` | `string` | Client secret (password) for an existing service principal, which will be granted access to the new vault. __WARNING: This secret will be shared with all users who bind to the vault!__ | Y | |
| `location` | `string` | The Azure region in which to provision applicable resources. | Required _unless_ an administrator has configured the broker itself with a default location. | The broker's default location, if configured. |
| `objectid` | `string` | Object ID for an existing service principal, which will be granted access to the new vault. | Y | |
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
| `vaultUri` | `string` | Fully qualified URI for connecting to the vault. |
| `clientId` | `string` | Service principal client ID (username) to use when connecting to the vault. |
| `clientSecret` | `string` | Service principal client secret (password) to use when connecting to the vault. |

##### Unbind

Does nothing.
  
##### Deprovision

Deletes the Key Vault.
