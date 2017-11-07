# Azure Redis Cache

[Azure Redis Cache](https://azure.microsoft.com/en-us/services/cache/) is based on the popular open-source Redis cache. It gives you access to a secure, dedicated Redis cache, managed by Microsoft and accessible from any application within Azure. This broker currently publishes a single service and plan for provisioning Azure Redis Cache.

## Services & Plans

### azure-rediscache

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
| `location` | `string` | The Azure region in which to provision applicable resources. | Y | |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | N | A new resource group will be created with a UUID as its name. |
| `tags` | `object` | Tags to be applied to new resources, specified as key/value pairs. | N | |
  
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
