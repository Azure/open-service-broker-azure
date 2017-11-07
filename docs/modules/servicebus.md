# Azure Service Bus

[Azure Service Bus](https://azure.microsoft.com/en-us/services/service-bus/) keep apps and devices connected across private and public clouds. This broker currently publishes a single service and plan for provisioning Azure Service Bus Service.

## Services & Plans

### azure-servicebus

| Plan Name | Description |
|-----------|-------------|
| `basic` | Basic Tier, Shared Capacity |
| `standard` | Standard Tier, Shared Capacity, Topics, 12.5M Messaging Operations/Month, Variable Pricing |
| `premium` | Premium Tier, Dedicated Capacity, Recommended For Production Workloads, Fixed Pricing |

#### Behaviors

##### Provision
  
Provisions a new Service Bus namespace. The new namespace will be named using new UUIDs.

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
| `connectionString` | `string` | Connection string. |
| `primaryKey` | `string` | Secret key (password). |

##### Unbind

Does nothing.
  
##### Deprovision

Deletes the Service Bus namespace.
