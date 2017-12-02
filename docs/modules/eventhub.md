# Azure Event Hubs

[Azure Event Hubs](https://azure.microsoft.com/en-us/services/event-hubs/) is a hyper-scale telemetry ingestion service that collects, transforms, and stores millions of events. As a distributed streaming platform, it gives you low latency and configurable time retention, which enables you to ingress massive amounts of telemetry into the cloud and read the data from multiple applications using publish-subscribe semantics.

## Services & Plans

### azure-eventhub

| Plan Name | Description |
|-----------|-------------|
| `basic` | Basic Tier, 1 Consumer group, 100 Brokered connections |
| `standard` | Standard Tier, 20 Consumer groups, 1000 Brokered connections, Additional Storage, Publisher Policies |

#### Behaviors

##### Provision
  
Provisions a new Event Hubs namespace and a new hub within it. The new namespace and hub will be named using new UUIDs.

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

Deletes the Event Hubs namespace.
