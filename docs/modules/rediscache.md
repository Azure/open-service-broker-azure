# [Azure Redis Cache](https://azure.microsoft.com/en-us/services/cache/)

_Note: This module is EXPERIMENTAL and is not included in the General Availability release of Open Service Broker for Azure. It will be added in a future OSBA release._

## Services & Plans

### Service: azure-rediscache

| Plan Name | Description |
|-----------|-------------|
| `basic` | Basic Tier, 250MB Cache |
| `standard` | Standard Tier, 1GB Cache |
| `standard` | Standard Tier, 1GB Cache |
| `premium` | Premium Tier, 6GB Cache |

#### Behaviors

##### Provision
  
Provisions a new Redis cache.

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
| `host` | `string` | The fully-qualified address of the Redis cache. |
| `port` | `int` | The port number to connect to on the Redis cache. |
| `password` | `string` | The password for the Redis cache. |

##### Unbind

Does nothing.
  
##### Deprovision

Deletes the Redis cache.
