# Azure SQL Database

[Azure SQL Database](https://azure.microsoft.com/en-us/documentation/articles/sql-database-technical-overview/) is a relational database service in the cloud based on the market-leading Microsoft SQL Server engine, with mission-critical capabilities.

## Services & Plans

### azure-sqldb

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
| `location` | `string` | The Azure region in which to provision applicable resources. | Y | |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | N | A new resource group will be created with a UUID as its name. |
| `server` | `string` | An optional reference to an existing server that has been pre-provsioned by a cluster admin, who has also pre-configured the broker with corresponding configuration for connecting to and administering that server. | N | |
| `tags` | `object` | Tags to be applied to new resources, specified as key/value pairs. | N | |
  
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
