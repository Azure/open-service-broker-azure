# [Azure SQL Database](https://azure.microsoft.com/en-us/services/sql-database/)

|![](https://upload.wikimedia.org/wikipedia/commons/thumb/1/17/Warning.svg/50px-Warning.svg.png) | This module is EXPERIMENTAL. It is under heavy development and remains subject to the possibility of breaking changes. |
|---|---|

## Services & Plans

### Service: azure-sqldb

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
  
By default, provisions a new SQL Server and a new database upon that server. The new database will be named randomly. If provisioning parameters include a reference to an existing server, provisioning a new server will be forgone and the new database will be provisioned upon the existing server. This option requires the server to have been pre-provsioned by a cluster admin, who has also pre-configured the broker with corresponding configuration for connecting to and administering that server.

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

If provisioning created a new SQL Server, deletes that SQL Server. If provisioning only created a new database on an existing SQL server, that database is dropped.
