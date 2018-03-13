# [Azure CosmosDB](https://azure.microsoft.com/en-us/services/cosmos-db/)

|![](https://upload.wikimedia.org/wikipedia/commons/thumb/1/17/Warning.svg/50px-Warning.svg.png) | This module is EXPERIMENTAL. It is under heavy development and remains subject to the possibility of breaking changes. |
|---|---|

## Services & Plans

### Service: azure-cosmos-db

| Plan Name | Description |
|-----------|-------------|
| `cosmos-db-sql-api` | Creates a CosmsDB Database Account configured to use SQL API |

#### Behaviors

##### Provision

Provisions a new CosmosDB database that can be accessed through any of the SQL API. The new database account is named using a new UUID.

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
| `uri` | `string` | The fully-qualified address and port of the CosmosDB database. ||
| `primarykey` | `string` | A secret key used for connecting to the CosmosDB database. |
| `primaryconnectionstring` | `string` | The full connection string, which includes the URI and primary key. |

##### Unbind

Does nothing.
  
##### Deprovision

Deletes the CosmosDB database.

### Service: azure-cosmos-mongo-db

| Plan Name | Description |
|-----------|-------------|
| `mongo-db` | MongoDB on Azure provided by CosmosDB |

#### Behaviors

##### Provision

Provisions a new CosmosDB database that can be accessed through the MongoDB API.

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
| `host` | `string` | The fully-qualified address of the CosmosDB database. |
| `port` | `int` | The port number to connect to on the CosmosDB database. |
| `username` | `string` | The name of the database user. |
| `password` | `string` | The password for the database user. |
| `connectionstring` | `string` | The full connection string, which includes the host, port, username, and password. |

##### Unbind

Does nothing.
  
##### Deprovision

Deletes the CosmosDB database.
