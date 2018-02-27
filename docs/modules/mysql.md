# [Azure Database for MySQL](https://azure.microsoft.com/en-us/services/mysql/)

|![](https://upload.wikimedia.org/wikipedia/commons/thumb/1/17/Warning.svg/50px-Warning.svg.png) | This module is EXPERIMENTAL. It is under heavy development and remains subject to the possibility of breaking changes. |
|---|---|

Open Service Broker for Azure contains three Azure Database for MySQL services. These services enable you to select the most appropriate provision scenario for your needs. These services are:

* azure-mysql - Provision both an Azure Database for MyQSL Database Management System (DBMS) and a database.
* azure-mysql-dbms-only - Provision only an Azure Database for MyQSL DBMS. This can be used to provision multiple databases at a later time.
* azure-mysql-db-only - Provision a new database only upon a previously provisioned DBMS.

The `azure-mysql` service allows you to provision both a DBMS and a database. This service is ready to use upon successful provisioing. You can not provision additional databases onto an instance provisioned through this service. The `azure-mysql-dbms-only` and `azure-mysql-db-only` services, on the other hand, can be combined to provison multiple databases on a single DBMS.  For more information on each service, refer to the descriptions below.

## Services & Plans

### Service: azure-mysql

| Plan Name | Description |
|-----------|-------------|
| `basic50` | Basic Tier, 50 DTUs |
| `basic100` | Basic Tier, 100 DTUs |

#### Behaviors

##### Provision

Provisions a new MySQL DBMS and a new database upon it. The new database will be named randomly.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Required _unless_ an administrator has configured the broker itself with a default location. | The broker's default location, if configured. |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | N | If an administrator has configured the broker itself with a default resource group and nonde is specified, that default will be applied, otherwise, a new resource group will be created with a UUID as its name. |
| `sslEnforcement` | `string` | Specifies whether the server requires the use of TLS when connecting. Valid valued are `""` (unspecified), `enabled`, or `disabled`. | N | `""`. Left unspecified, SSL _will_ be enforced. |
| `firewallRules`  | `array` | Specifies the firewall rules to apply to the server. Definition follows. | N | `[]` Left unspecified, Firewall will default to only Azure IPs. If rules are provided, they must have valid values. |
| `firewallRules[n].name` | `string` | Specifies the name of the generated firewall rule |Y | |
| `firewallRules[n].startIPAddress` | `string` | Specifies the start of the IP range allowed by this firewall rule | Y | |
| `firewallRules[n].endIPAddress` | `string` | Specifies the end of the IP range allowed by this firewall rule | Y | |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |

##### Bind

Creates a new user on the MySQL DBMS. The new user will be named randomly and
will be granted a wide array of permissions on the database.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and credentials:

| Field Name | Type | Description |
|------------|------|-------------|
| `host` | `string` | The fully-qualified address of the MySQL DBMS. |
| `port` | `int` | The port number to connect to on the MySQL DBMS. |
| `database` | `string` | The name of the database. |
| `username` | `string` | The name of the database user (in the form username@host). |
| `password` | `string` | The password for the database user. |

##### Unbind

Drops the applicable user from the MySQL DBMS.

##### Deprovision

Deletes the MySQL DBMS and database.

### Service: azure-mysql-dbms-only

| Plan Name | Description |
|-----------|-------------|
| `basic50` | Basic Tier, 50 DTUs |
| `basic100` | Basic Tier, 100 DTUs |

#### Behaviors

##### Provision

Provisions a new MySQL DBMS only. Databases can be created through subsequent provision reqquests using the `azure-mysql-database-only` service.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Required _unless_ an administrator has configured the broker itself with a default location. | The broker's default location, if configured. |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | N | If an administrator has configured the broker itself with a default resource group and nonde is specified, that default will be applied, otherwise, a new resource group will be created with a UUID as its name. |
| `alias` | `string` | Y | Specifies an alias that can be used by later provision actions to create databases on this DBMS. |
| `sslEnforcement` | `string` | Specifies whether the server requires the use of TLS when connecting. Valid valued are `""` (unspecified), `enabled`, or `disabled`. | N | `""`. Left unspecified, SSL _will_ be enforced. |
| `firewallRules`  | `array` | Specifies the firewall rules to apply to the server. Definition follows. | N | `[]` Left unspecified, Firewall will default to only Azure IPs. If rules are provided, they must have valid values. |
| `firewallRules[n].name` | `string` | Specifies the name of the generated firewall rule |Y | |
| `firewallRules[n].startIPAddress` | `string` | Specifies the start of the IP range allowed by this firewall rule | Y | |
| `firewallRules[n].endIPAddress` | `string` | Specifies the end of the IP range allowed by this firewall rule | Y | |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |

##### Bind

This service is not bindable.

###### Binding Parameters

This service is not bindable.

###### Credentials

This service is not bindable.

##### Unbind

This service is not bindable.

##### Deprovision

Deprovision will delete the MySQL DBMS. If any databases have been provisioned on this DBMS, deprovisioning will be deferred until all databases have been deprovisioned.

### Service: azure-mysqldb-db-only

| Plan Name | Description |
|-----------|-------------|
| `mysql-db-only` | New database on existing MySQL DBMS |

#### Behaviors

##### Provision

Provisions a new database upon a previously provisioned DBMS. The new database will be named randomly.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `parentAlias` | `string` | Specifies the alias of the DBMS upon which the database should be provisioned. | Y | |

##### Bind

Creates a new user on the MySQL DBMS. The new user will be named randomly and
will be granted a wide array of permissions on the database.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and credentials:

| Field Name | Type | Description |
|------------|------|-------------|
| `host` | `string` | The fully-qualified address of the MySQL DBMS. |
| `port` | `int` | The port number to connect to on the MySQL DBMS. |
| `database` | `string` | The name of the database. |
| `username` | `string` | The name of the database user (in the form username@host). |
| `password` | `string` | The password for the database user. |

##### Unbind

Drops the applicable user from the MySQL DBMS.

##### Deprovision

Deletes the database from the MySQL DBMS. The DBMS itself is not deprovisioned.