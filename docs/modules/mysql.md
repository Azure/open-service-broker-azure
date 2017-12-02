# Azure Database for MySQL


[Azure Database for MySQL](https://azure.microsoft.com/en-us/services/mysql) is a relational database service in the Microsoft cloud based on MySQL Community Edition database engine. Azure Database for MySQL delivers:

* Predictable performance at multiple service levels
* Dynamic scalability with no application downtime
* Built-in high availability
* Data protection

## Services & Plans

### azure-mysqldb

| Plan Name | Description |
|-----------|-------------|
| `basic50` | Basic Tier, 50 DTUs |
| `basic100` | Basic Tier, 100 DTUs |

#### Behaviors

##### Provision
  
Provisions a new MySQL server and a new database upon that server. The new database will be named randomly.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Y | |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | N | A new resource group will be created with a UUID as its name. |
| `sslEnforcement` | `string` | Specifies whether the server requires the use of TLS when connecting. Valid valued are `""` (unspecified), `enabled`, or `disabled`. | N | `""`. Left unspecified, SSL _will_ be enforced. |
| `tags` | `object` | Tags to be applied to new resources, specified as key/value pairs. | N | |
  
##### Bind
  
Creates a new user on the MySQL server. The new user will be named randomly and will be granted a wide array of permissions on the database.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and credentials:

| Field Name | Type | Description |
|------------|------|-------------|
| `host` | `string` | The fully-qualified address of the MySQL server. |
| `port` | `int` | The port number to connect to on the MySQL server. |
| `database` | `string` | The name of the database. |
| `username` | `string` | The name of the database user (in the form username@host). |
| `password` | `string` | The password for the database user. |

##### Unbind

Drops the applicable user from the MySQL server.
  
##### Deprovision

Deletes the MySQL server.
