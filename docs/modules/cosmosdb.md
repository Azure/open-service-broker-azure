# [Azure CosmosDB](https://azure.microsoft.com/en-us/services/cosmos-db/)

_Note: This module is EXPERIMENTAL and future releases may break the API._

## Services & Plans

### Service: azure-cosmosdb-sql-account

| Plan Name | Description |
|-----------|-------------|
| `account` | Database Account configured to use SQL API |

#### Behaviors

##### Provision

Provisions a new CosmosDB database account that can be accessed through any of the SQL API. The new database account is named using a new UUID.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Required _unless_ an administrator has configured the broker itself with a default location. | The broker's default location, if configured. |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | N | If an administrator has configured the broker itself with a default resource group and nonde is specified, that default will be applied, otherwise, a new resource group will be created with a UUID as its name. |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |
| `ipFilters` | `object` | IP Range Filter to be applied to new CosmosDB account | N | No Filters Applied |
| `ipFilters.allowAccessFromAzure` | `string` | Specifies if Azure Services should be able to access the CosmosDB account.Valid valued are `""` (unspecified), `enabled`, or `disabled`. | N | If left unspecified, defaults to enabled. |
| `ipFilters.allowAccessFromPortal` | `string` | Specifies if the Azure Portal should be able to access the CosmosDB account. If `allowAccessFromAzure` is set to enabled, this value is ignored. Valid valued are `""` (unspecified), `enabled`, or `disabled`. | N | If left unspecified, defaults to enabled. |
| `ipFilters.filters` | `array` | Values to include in IP Filter. Can be IP Address or CIDR range. | N | If not specified, no additional values will be included in filters. |

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

### Service: azure-cosmosdb-mongo-account

| Plan Name | Description |
|-----------|-------------|
| `account` | MongoDB on Azure provided by CosmosDB |

#### Behaviors

##### Provision

Provisions a new CosmosDB database account that can be accessed through the MongoDB API. The new database account is named using a new UUID.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Required _unless_ an administrator has configured the broker itself with a default location. | The broker's default location, if configured. |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | N | If an administrator has configured the broker itself with a default resource group and nonde is specified, that default will be applied, otherwise, a new resource group will be created with a UUID as its name. |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |
| `ipFilters` | `object` | IP Range Filter to be applied to new CosmosDB account | N | No Filters Applied |
| `ipFilters.allowAccessFromAzure` | `string` | Specifies if Azure Services should be able to access the CosmosDB account.Valid valued are `""` (unspecified), `enabled`, or `disabled`. | N | If left unspecified, defaults to enabled. |
| `ipFilters.allowAccessFromPortal` | `string` | Specifies if the Azure Portal should be able to access the CosmosDB account. If `allowAccessFromAzure` is set to enabled, this value is ignored. Valid valued are `""` (unspecified), `enabled`, or `disabled`. | N | If left unspecified, defaults to enabled. |
| `ipFilters.filters` | `array` | Values to include in IP Filter. Can be IP Address or CIDR range. | N | If not specified, no additional values will be included in filters. |

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

### Service: azure-cosmosdb-graph-account

| Plan Name | Description |
|-----------|-------------|
| `account` | Database Account configured to use Graph (Gremlin) API |

#### Behaviors

##### Provision

Provisions a new CosmosDB database account that can be accessed through any of the Graph (Gremlin) API. The new database account is named using a new UUID.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Required _unless_ an administrator has configured the broker itself with a default location. | The broker's default location, if configured. |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | N | If an administrator has configured the broker itself with a default resource group and nonde is specified, that default will be applied, otherwise, a new resource group will be created with a UUID as its name. |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |
| `ipFilters` | `object` | IP Range Filter to be applied to new CosmosDB account | N | No Filters Applied |
| `ipFilters.allowAccessFromAzure` | `string` | Specifies if Azure Services should be able to access the CosmosDB account.Valid valued are `""` (unspecified), `enabled`, or `disabled`. | N | If left unspecified, defaults to enabled. |
| `ipFilters.allowAccessFromPortal` | `string` | Specifies if the Azure Portal should be able to access the CosmosDB account. If `allowAccessFromAzure` is set to enabled, this value is ignored. Valid valued are `""` (unspecified), `enabled`, or `disabled`. | N | If left unspecified, defaults to enabled. |
| `ipFilters.filters` | `array` | Values to include in IP Filter. Can be IP Address or CIDR range. | N | If not specified, no additional values will be included in filters. |

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

### Service: azure-cosmosdb-table-account

| Plan Name | Description |
|-----------|-------------|
| `account` | Database Account configured to use Table API |

#### Behaviors

##### Provision

Provisions a new CosmosDB database account that can be accessed through any of the Azure Table API. The new database account is named using a new UUID.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Required _unless_ an administrator has configured the broker itself with a default location. | The broker's default location, if configured. |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | N | If an administrator has configured the broker itself with a default resource group and nonde is specified, that default will be applied, otherwise, a new resource group will be created with a UUID as its name. |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |
| `ipFilters` | `object` | IP Range Filter to be applied to new CosmosDB account | N | No Filters Applied |
| `ipFilters.allowAccessFromAzure` | `string` | Specifies if Azure Services should be able to access the CosmosDB account.Valid valued are `""` (unspecified), `enabled`, or `disabled`. | N | If left unspecified, defaults to enabled. |
| `ipFilters.allowAccessFromPortal` | `string` | Specifies if the Azure Portal should be able to access the CosmosDB account. If `allowAccessFromAzure` is set to enabled, this value is ignored. Valid valued are `""` (unspecified), `enabled`, or `disabled`. | N | If left unspecified, defaults to enabled. |
| `ipFilters.filters` | `array` | Values to include in IP Filter. Can be IP Address or CIDR range. | N | If not specified, no additional values will be included in filters. |

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