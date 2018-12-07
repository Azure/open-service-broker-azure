# [Azure Service Bus](https://azure.microsoft.com/en-us/services/service-bus/)

_Note: This module is EXPERIMENTAL. To enable this module, you must run Open Service Broker for Azure with modules.minStability set to `experimental`_

_This module involves the Parent-Child Model concept in OSBA, please refer to the [Parent-Child Model doc](../parent-child-model-for-multiple-layers-services.md)._

## Services & Plans

### Service: azure-servicebus-namespace

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
| `location` | `string` | The Azure region in which to provision applicable resources. | Y |  |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | Y |  |
| `alias` | `string` | Specifies an alias that can be used by later provision actions to create queues/topics in this namespace. | Y |  |
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



### Service: azure-servicebus-queue

| Plan Name | Description                     |
| --------- | ------------------------------- |
| `queue`   | New queue in existing namespace |

#### Behaviors

##### Provision

Provisions a new queue in an existing namespace. 

###### Provisioning Parameters

| Parameter Name      | Type     | Description                                                  | Required | Default Value                                                |
| ------------------- | -------- | ------------------------------------------------------------ | -------- | ------------------------------------------------------------ |
| `parentAlias`       | `string` | Specifies the alias of the namespace in which the  queue should be provisioned. | Y        |                                                              |
| `queueName`         | `string` | The name of the queue.                                       | N        | If not provided, a random name will be generated as the queue name. |
| `maxQueueSize`      | `int`    | The maximum size of the queue in megabytes, which is the size of memory allocated for the queue. | N        | 1024                                                         |
| `messageTimeToLive` | `string` | ISO 8601 default message timespan to live value. This is the duration after which the message expires, starting from when the message is sent to Service Bus. This is the default value used when TimeToLive is not set on a message itself. For example, `PT276H13M14S` sets the message to expire in 11 day 12 hour 13 minute 14 seconds. | N        | "PT336H"                                                     |
| `lockDuration`      | `string` | ISO 8601 timespan duration of a peek-lock; that is, the amount of time that the message is locked for other receivers. The lock duration time window can range from 5 seconds to 5 minutes. For example, `PT2M30S` sets the lock duration time to 2 minutes 30 seconds. | N        | "PT30S"                                                      |

##### Bind

Returns a copy of one shared set of credentials.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and shared credentials:

| Field Name         | Type     | Description            |
| ------------------ | -------- | ---------------------- |
| `connectionString` | `string` | Connection string.     |
| `primaryKey`       | `string` | Secret key (password). |
| `queueURL`         | `string` | Queue URL              |

##### Unbind

Does nothing.

##### Deprovision

Deletes the Service Bus queue.



### Service: azure-servicebus-topic

| Plan Name | Description                     |
| --------- | ------------------------------- |
| `topic`   | New topic in existing namespace |

#### Behaviors

##### Provision

Provisions a new topic in an existing namespace. 

###### Provisioning Parameters

| Parameter Name      | Type     | Description                                                  | Required | Default Value                                                |
| ------------------- | -------- | ------------------------------------------------------------ | -------- | ------------------------------------------------------------ |
| `parentAlias`       | `string` | Specifies the alias of the namespace in which the  topic should be provisioned. **Note**: the parent must be a service-bus-namespace instance with `standard` or `premium` plan. | Y        |                                                              |
| `topicName`         | `string` | The name of the topic                                        | N        | If not provided, a random name will be generated as the topic name. |
| `maxTopicSize`      | `int`    | The maximum size of the topic in megabytes, which is the size of memory allocated for the topic. | N        | 1024                                                         |
| `messageTimeToLive` | `string` | ISO 8601 default message timespan to live value. This is the duration after which the message expires, starting from when the message is sent to Service Bus. This is the default value used when TimeToLive is not set on a message itself. For example, `PT276H13M14S` sets the message to expire in 11 days 12 hours 13 minutes 14 seconds. | N        | "PT336H"                                                     |

##### Bind

Returns a copy of one shared set of credentials.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and shared credentials:

| Field Name         | Type     | Description            |
| ------------------ | -------- | ---------------------- |
| `connectionString` | `string` | Connection string.     |
| `primaryKey`       | `string` | Secret key (password). |
| `topicURL`         | `string` | Topic URL              |

##### Unbind

Does nothing.

##### Deprovision

Deletes the Service Bus topic.