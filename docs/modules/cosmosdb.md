# Azure CosmosDB

[Azure CosmosDB](https://azure.microsoft.com/en-us/services/cosmos-db/) is a globally distributed database service designed to enable you to elastically and independently scale throughput and storage across any number of geographical regions with a comprehensive SLA. You can develop document, key/value, or graph databases with Cosmos DB using a series of popular APIs and programming models. Learn how to use Cosmos DB with our quickstarts, tutorials, and samples.

## Services & Plans

### azure-cosmos-document-db

| Plan Name | Description |
|-----------|-------------|
| `document-db` | Azure DocumentDB provided by CosmosDB and accessible via SQL (DocumentDB), Gremlin (Graph), and Table (Key-Value) APIs |

#### Behaviors

##### Provision
  
Provisions a new CosmosDB database that can be accessed through any of the SQL (DocumentDB), Gremlin (Graph), and Table (Key-Value) APIs. The new database is named using a new UUID.

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
| `uri` | `string` | The fully-qualified address and port of the CosmosDB database. ||
| `primarykey` | `string` | A secret key used for connecting to the CosmosDB database. |
| `primaryconnectionstring` | `string` | The full connection string, which includes the URI and primary key. |

##### Unbind

Does nothing.
  
##### Deprovision

Deletes the CosmosDB database.

### azure-cosmos-mongo-db

| Plan Name | Description |
|-----------|-------------|
| `azure-cosmos-mongo-db` | MongoDB on Azure provided by CosmosDB |

#### Behaviors

##### Provision
  
Provisions a new CosmosDB database that can be accessed through the MongoDB API.

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
| `host` | `string` | The fully-qualified address of the CosmosDB database. |
| `port` | `int` | The port number to connect to on the CosmosDB database. |
| `username` | `string` | The name of the database user. |
| `password` | `string` | The password for the database user. |
| `connectionstring` | `string` | The full connection string, which includes the host, port, username, and password. |

##### Unbind

Does nothing.
  
##### Deprovision

Deletes the CosmosDB database.
