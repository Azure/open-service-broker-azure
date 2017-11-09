# Azure Key Vault

## Services & Plans

### azure-keyvault

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
| `location` | `string` | The Azure region in which to provision applicable resources. | Y | |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | N | A new resource group will be created with a UUID as its name. |
| `objectid` | `string` | Object ID for an existing service principal, which will be granted access to the new vault. | Y | |
| `clientId` | `string` | Client ID (username) for an existing service principal, which will be granted access to the new vault.| Y | |
| `clientSecret` | `string` | Client secret (password) for an existing service principal, which will be granted access to the new vault. __WARNING: This secret will be shared with all users who bind to the vault!__ | Y | |
| `tags` | `object` | Tags to be applied to new resources, specified as key/value pairs. | N | |
  
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
