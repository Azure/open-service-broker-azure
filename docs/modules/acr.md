# [Azure Container Registry](https://azure.microsoft.com/en-us/services/container-registry/)

|![](https://upload.wikimedia.org/wikipedia/commons/thumb/1/17/Warning.svg/50px-Warning.svg.png) | This module is EXPERIMENTAL. It is under heavy development and remains subject to the possibility of breaking changes. |
|---|---|

## Services & Plans

### Service: azure-containerregsitry

| Plan Name | Description |
|-----------|-------------|
| `Basic` | Basic Tier |
| `Standard` | Standard Tier |
| `Premium` | Premium Tier |

#### Behaviors

##### Provision
  
Provisions an Azure Container Registry.

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

Binding returns the following connection details and credentials:

| Field Name | Type | Description |
|------------|------|-------------|
| `registryName` | `string` | The name of the Azure Container Registry. |

##### Unbind

Drops the applicable role (user) from the Container Registry.
  
##### Deprovision

Deletes the Container Registry.
