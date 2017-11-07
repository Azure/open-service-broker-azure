# Azure Database for PostgreSQL

[Azure Database for PostgreSQL](https://azure.microsoft.com/en-us/services/postgresql) is a relational database service based on the open source Postgres database engine. It is a fully managed database as a service offering capable of handling mission-critical workloads with predictable performance, security, high availability, and dynamic scalability.  Develop applications with Azure Database for PostgreSQL leveraging the open source tools and platform of your choice.

## Services & Plans

### azure-postgresqldb

| Plan Name | Description |
|-----------|-------------|
| `basic50` | Basic Tier, 50 DTUs |
| `basic100` | Basic Tier, 100 DTUs |

#### Behaviors

##### Provision
  
Provisions a new PostgreSQL server and a new database upon that server. The new database will be named randomly and will be owned by a role (group) of the same name.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Y | |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | N | A new resource group will be created with a UUID as its name. |
| `sslEnforcement` | `string` | Specifies whether the server requires the use of TLS when connecting. Valid valued are `""` (unspecified), `enabled`, or `disabled`. | N | `""`. Left unspecified, SSL _will_ be enforced. |
| `tags` | `object` | Tags to be applied to new resources, specified as key/value pairs. | N | |
  
##### Bind
  
Creates a new role (user) on the PostgreSQL server. The new role will be named randomly and added to the  role (group) that owns the database.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and credentials:

| Field Name | Type | Description |
|------------|------|-------------|
| `host` | `string` | The fully-qualified address of the PostgreSQL server. |
| `port` | `int` | The port number to connect to on the PostgreSQL server. |
| `database` | `string` | The name of the database. |
| `username` | `string` | The name of the database user (in the form username@host). |
| `password` | `string` | The password for the database user. |

##### Unbind

Drops the applicable role (user) from the PostgreSQL server.
  
##### Deprovision

Deletes the PostgreSQL server.
