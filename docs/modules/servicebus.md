# [Azure Service Bus](https://azure.microsoft.com/en-us/services/service-bus/)

_Note: This module is EXPERIMENTAL and future releases may break the API._

## Services & Plans

### Service: azure-servicebus

| Plan Name | Description |
|-----------|-------------|
| `basic` | Basic Tier, Shared Capacity |
| `standard` | Standard Tier, Shared Capacity, Topics, 12.5M Messaging Operations/Month, Variable Pricing |
| `premium` | Premium Tier, Dedicated Capacity, Recommended For Production Workloads, Fixed Pricing |

#### Behaviors

##### Provision
  
Provisions a new Service Bus namespace. The new namespace will be named using
new UUIDs.

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
| `connectionString` | `string` | Connection string. |
| `primaryKey` | `string` | Secret key (password). |

##### Unbind

Does nothing.
  
##### Deprovision

Deletes the Service Bus namespace.
