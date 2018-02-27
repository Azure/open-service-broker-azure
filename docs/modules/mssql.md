# [Azure SQL Database](https://azure.microsoft.com/en-us/services/sql-database/)

|![](https://upload.wikimedia.org/wikipedia/commons/thumb/1/17/Warning.svg/50px-Warning.svg.png) | This module is EXPERIMENTAL. It is under heavy development and remains subject to the possibility of breaking changes. |
|---|---|

Open Service Broker for Azure contains three Azure SQL Database services. These services enable you to select the most appropriate provision scenario for your needs. These services are:

* azure-sql - Provision both a SQL Server DBMS and a database
* azure-sql-dbms-only - Provision only a SQL Server Database Management System (DBMS). This can be used to provision multiple databases at a later time.
* azure-sql-db-only - Provision a new database only upon a previously provisioned DBMS.

The `azure-sql` service allows you to provision both a DBMS and a database. This service is ready to use upon successful provisioing. You can not provision additional databases onto an instance provisioned through this service. The `azure-sql-dbms-only` and `azure-sql-db-only` services, on the other hand, can be combined to provison multiple databases on a single DBMS.  For more information on each service, refer to the descriptions below.

## Services & Plans

### Service: azure-sql

| Plan Name | Description |
|-----------|-------------|
| `basic` | Basic Tier, 5 DTUs, 2GB, 7 days point-in-time restore |
| `standard-s0` | Standard Tier, 10 DTUs, 250GB, 35 days point-in-time restore |
| `standard-s1` | StandardS1 Tier, 20 DTUs, 250GB, 35 days point-in-time restore |
| `standard-s2` | StandardS2 Tier, 50 DTUs, 250GB, 35 days point-in-time restore |
| `standard-s3` | StandardS3 Tier, 100 DTUs, 250GB, 35 days point-in-time restore |
| `premium-p1` | PremiumP1 Tier, 125 DTUs, 500GB, 35 days point-in-time restore |
| `premium-p2` | PremiumP2 Tier, 250 DTUs, 500GB, 35 days point-in-time restore |
| `premium-p4` | PremiumP4 Tier, 500 DTUs, 500GB, 35 days point-in-time restore |
| `premium-p6` | PremiumP6 Tier, 1000 DTUs, 500GB, 35 days point-in-time restore |
| `premium-p11` | PremiumP11 Tier, 1750 DTUs, 1024GB, 35 days point-in-time restore |
| `data-warehouse-100` | DataWarehouse100 Tier, 100 DWUs, 1024GB |
| `data-warehouse-1200` | DataWarehouse1200 Tier, 1200 DWUs, 1024GB |

#### Behaviors

##### Provision

Provisions a new SQL Server and a new database upon that server. The new dbms and database will be named randomly. 

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Required _unless_ an administrator has configured the broker itself with a default location. | The broker's default location, if configured. |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | N | If an administrator has configured the broker itself with a default resource group and nonde is specified, that default will be applied, otherwise, a new resource group will be created with a UUID as its name. |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |
| `firewallRules`  | `array` | Specifies the firewall rules to apply to the server. Definition follows. | N | `[]` Left unspecified, Firewall will default to only Azure IPs. If rules are provided, they must have valid values. |
| `firewallRules[n].name` | `string` | Specifies the name of the generated firewall rule |Y | |
| `firewallRules[n].startIPAddress` | `string` | Specifies the start of the IP range allowed by this firewall rule | Y | |
| `firewallRules[n].endIPAddress` | `string` | Specifies the end of the IP range allowed by this firewall rule | Y | |

##### Bind

Creates a new user on the SQL Server. The new user will be named randomly and granted permission to log into and administer the database.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and credentials:

| Field Name | Type | Description |
|------------|------|-------------|
| `host` | `string` | The fully-qualified address of the SQL Server. |
| `port` | `int` | The port number to connect to on the SQL Server. |
| `database` | `string` | The name of the database. |
| `username` | `string` | The name of the database user (in the form username@host). |
| `password` | `string` | The password for the database user. |

##### Unbind

Drops the applicable user from the SQL Server.

##### Deprovision

Deletes both the database and the SQL Server instance.

### Service: azure-sql-dbms-only

| Plan Name | Description |
|-----------|-------------|
| `sql-dbms-only` | Azure SQL Server DBMS Only|

#### Behaviors

##### Provision

Provisions a new SQL Server DBMS only. Databases can be created through subsequent provision reqquests using the `azure-sql-db-only` service.
###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Required _unless_ an administrator has configured the broker itself with a default location. | The broker's default location, if configured. |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | N | If an administrator has configured the broker itself with a default resource group and nonde is specified, that default will be applied, otherwise, a new resource group will be created with a UUID as its name. |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |
| `alias` | `string` | Y | Specifies an alias that can be used by later provision actions to create databases on this DBMS. |
| `firewallRules`  | `array` | Specifies the firewall rules to apply to the server. Definition follows. | N | `[]` Left unspecified, Firewall will default to only Azure IPs. If rules are provided, they must have valid values. |
| `firewallRules[n].name` | `string` | Specifies the name of the generated firewall rule |Y | |
| `firewallRules[n].startIPAddress` | `string` | Specifies the start of the IP range allowed by this firewall rule | Y | |
| `firewallRules[n].endIPAddress` | `string` | Specifies the end of the IP range allowed by this firewall rule | Y | |

##### Bind

This service is not bindable.

###### Binding Parameters

This service is not bindable.

###### Credentials

This service is not bindable.

##### Unbind

This service is not bindable.

##### Deprovision

Deprovision will delete the SQL Server DBMS. If any databases have been provisioned on this DBMS, deprovisioning will be deferred until all databases have been deprovisioned.

### Service: azure-sql-db-only

| Plan Name | Description |
|-----------|-------------|
| `basic` | Basic Tier, 5 DTUs, 2GB, 7 days point-in-time restore |
| `standard-s0` | Standard Tier, 10 DTUs, 250GB, 35 days point-in-time restore |
| `standard-s1` | StandardS1 Tier, 20 DTUs, 250GB, 35 days point-in-time restore |
| `standard-s2` | StandardS2 Tier, 50 DTUs, 250GB, 35 days point-in-time restore |
| `standard-s3` | StandardS3 Tier, 100 DTUs, 250GB, 35 days point-in-time restore |
| `premium-p1` | PremiumP1 Tier, 125 DTUs, 500GB, 35 days point-in-time restore |
| `premium-p2` | PremiumP2 Tier, 250 DTUs, 500GB, 35 days point-in-time restore |
| `premium-p4` | PremiumP4 Tier, 500 DTUs, 500GB, 35 days point-in-time restore |
| `premium-p6` | PremiumP6 Tier, 1000 DTUs, 500GB, 35 days point-in-time restore |
| `premium-p11` | PremiumP11 Tier, 1750 DTUs, 1024GB, 35 days point-in-time restore |
| `data-warehouse-100` | DataWarehouse100 Tier, 100 DWUs, 1024GB |
| `data-warehouse-1200` | DataWarehouse1200 Tier, 1200 DWUs, 1024GB |

#### Behaviors

##### Provision

Provisions a new database upon an existing server. The new database will be named randomly. If the DBMS does not yet exist, provision of the database will be deferred until the DBMS has been provisioned.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `parentAlias` | `string` | Specifies the alias of the DBMS upon which the database should be provisioned. | Y | |

##### Bind

Creates a new user on the SQL Server. The new user will be named randomly and granted permission to log into and administer the database.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and credentials:

| Field Name | Type | Description |
|------------|------|-------------|
| `host` | `string` | The fully-qualified address of the SQL Server. |
| `port` | `int` | The port number to connect to on the SQL Server. |
| `database` | `string` | The name of the database. |
| `username` | `string` | The name of the database user (in the form username@host). |
| `password` | `string` | The password for the database user. |

##### Unbind

Drops the applicable user from the SQL Server.

##### Deprovision

Deletes both the database from the SQL Server instance, but does not delete the DBMS.