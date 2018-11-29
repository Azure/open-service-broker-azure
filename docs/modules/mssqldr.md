# [Azure SQL Database Failover Group](https://docs.microsoft.com/en-us/azure/sql-database/sql-database-geo-replication-overview)

_Note: The services in this module aren't totally addressed.`_

_Note: This module is EXPERIMENTAL. To enable this module, you must run Open Service Broker for Azure with modules.minStability set to `experimental`_

Open Service Broker for Azure (OSBA) contains five Azure SQL Database Failover Group services. These services enable you to select the most appropriate provisioning scenario for your needs. These services are:

| Service Name | Description |
|--------------|-------------|
| `azure-sql-12-0-dr-dbms-pair-registered` | Register two existing servers as a service instance. |
| `azure-sql-12-0-dr-database-pair` | Provision two new databases which make up a failover group upon a previous DBMS pair. |
| `azure-sql-12-0-dr-database-pair-registered` | Register a failover group upon a previous DBMS pair as a service instance. |
| `azure-sql-12-0-dr-database-pair-from-existing` | Taking over an existing failover group (included the databases) upon a previous DBMS pair as a service instance. |
| `azure-sql-12-0-dr-database-pair-from-existing-primary` | Taking over an existing database upon the primary server of a previous DBMS pair, and extend it to a failover group deployment. |

All the services in this module require `ENABLE_DISASTER_RECOVERY_SERVICES` to be `true` in OSBA environment variables. Besides, `azure-sql-12-0-dr-database-pair-from-existing` and `azure-sql-12-0-dr-database-pair-from-existing-primary` require `ENABLE_MIGRATION_SERVICES` to be `true`. For more information on each service, refer to the descriptions below.

_This module involves the Parent-Child Model concept in OSBA, please refer to the [Parent-Child Model doc](../parent-child-model-for-multiple-layers-services.md)_.

## Services & Plans

### Service: azure-sql-12-0-dr-dbms-pair-registered

| Plan Name | Description |
|-----------|-------------|
| `dbms` | Azure SQL Server DBMS-Only |

#### Behaviors

##### Provision

Register a pair of SQL servers as a service instance: check the existence of these servers; check if the input administrator logins work. Databases with failover group can be created through subsequent provision requests using the `azure-sql-12-0-dr-database-pair` service.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `primaryResourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | Y | |
| `primaryLocation` | `string` | The Azure region in which to provision applicable resources. | Y | |
| `primaryServer` | `string` | The name of your existing primary server. | Y | |
| `primaryAdministratorLogin` | `string` | The administrator login of the primary server. | Y | |
| `primaryAdministratorLoginPassword` | `string` | The administrator login password of the primary server. | Y | |
| `secondaryResourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | Y | |
| `secondaryLocation` | `string` | The Azure region in which to provision applicable resources. | Y | |
| `secondaryServer` | `string` | The name of your existing secondary server. | Y | |
| `secondaryAdministratorLogin` | `string` | The administrator login of the secondary server. | Y | |
| `secondaryAdministratorLoginPassword` | `string` | The administrator login password of the secondary server. | Y | |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |
| `alias` | `string` | Specifies an alias that can be used by later provision actions to create database pairs on this DBMS pair. | Y | |

##### Bind

This service is not bindable.

##### Update

Updates broker-stored administrator login/password in case you reset them.

###### Updating Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `primaryAdministratorLogin` | `string` | The administrator login of the primary server. | N | |
| `primaryAdministratorLoginPassword` | `string` | The administrator login password of the primary server. | N | |
| `secondaryAdministratorLogin` | `string` | The administrator login of the secondary server. | N | |
| `secondaryAdministratorLoginPassword` | `string` | The administrator login password of the secondary server. | N | |

##### Unbind

This service is not bindable.

##### Deprovision

Do nothing as it is a registered instance. If any database pairs have been provisioned on this DBMS pair, deprovisioning will be deferred until all databases have been deprovisioned.

##### Examples

###### Kubernetes

To add.

###### Cloud Foundry

Using the `cf` cli, you can create the `dbms` plan of the `azure-sql-12-0-dr-dbms-pair-registered` service with the following command:

```console
cf create-service azure-sql-12-0-dr-dbms-pair-registered dbms serverpair -c '{
  "primaryResourceGroup":"osba",
  "primaryServer":"osbasql1",
  "primaryLocation":"eastus",
  "primaryAdministratorLogin":"username1",
  "primaryAdministratorLoginPassword":"password1",
  "secondaryResourceGroup":"osba",
  "secondaryServer":"osbasql2",
  "secondaryLocation":"westus",
  "secondaryAdministratorLogin":"username2",
  "secondaryAdministratorLoginPassword":"password2",
  "alias":"serverpair"
}'
```

###### cURL

To provision an instance using the broker directly, you must use the ID for both plan and service. Assuming your OSBA is running locally on port 8080 with the default username and password, you can provision the DBMS-only Azure SQL Database service with a cURL command similar to the following example. This example illustrates multiple firewall rules and provides an alias for later database provisioning:

```console
curl -X PUT \
  'http://localhost:8080/v2/service_instances/sql-dbms?accepts_incomplete=true' \
  -H 'authorization: Basic dXNlcm5hbWU6cGFzc3dvcmQ=' \
  -H 'content-type: application/json' \
  -H 'x-broker-api-version: 2.13' \
  -d '{
    "service_id" : "a7454e0e-be2c-46ac-b55f-8c4278117525",
    "plan_id" : "24f0f42e-1ab3-474e-a5ca-b943b2c48eee",
    "parameters" : {
      "primaryResourceGroup":"osba",
      "primaryServer":"osbasql1",
      "primaryLocation":"eastus",
      "primaryAdministratorLogin":"username1",
      "primaryAdministratorLoginPassword":"password1",
      "secondaryResourceGroup":"osba",
      "secondaryServer":"osbasql2",
      "secondaryLocation":"westus",
      "secondaryAdministratorLogin":"username2",
      "secondaryAdministratorLoginPassword":"password2",
      "alias":"serverpair"
    }
}'
```

### Service: azure-sql-12-0-dr-database-pair

| Plan Name | Description |
|-----------|-------------|
| `basic` | Basic Tier, 5 DTUs, 2GB, 7 days point-in-time restore |
| `standard` | Standard Tier, Up to 3000 DTUs, with 250GB storage, 35 days point-in-time restore |
| `premium` | Premium Tier, Up to 4000 DTUs, with 500GB storage, 35 days point-in-time restore |
| `general-purpose` | General Purpose Tier, Up to 80 vCores, Up to 440 GB Memory, Up to 1 TB storage, 7 days point-in-time restore |
| `business-critical` | Business Critical Tier, Up to 80 vCores, Up to 440 GB Memory, Up to 1 TB storage, Local SSD, 7 days point-in-time restore. Offers highest resilience to failures using several isolated replicas |

#### Behaviors

##### Provision

Provisions a pair of new databases upon both the primary server and the secondary server. And these two databases make up a failover group. If the DBMS pair does not yet exist, provision of the database will be deferred until the DBMS has been provisioned.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `parentAlias` | `string` | Specifies the alias of the DBMS upon which the database should be provisioned. | Y | |
| `database` | `string` | Specifies the name of the databases. | Y | |
| `failoverGroup` | `string` | Specifies the name of the failover group. | Y | |

Additional Provision Parameters for : standard plan

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `dtus` | `integer` | Specifies Database transaction units, which represent a bundled measure of compute, storage, and IO resources. Valid values are 10, 20, 50, 100, 200, 400, 800, 1600, 3000 | N | 10 |


Additional Provision Parameters for : premium plan

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `dtus` | `integer` | Specifies Database transaction units, which represent a bundled measure of compute, storage, and IO resources. Valid values are 125, 250, 500, 1000, 1750, 4000 | N | 125 |

Additional Provision Parameters for: general-purpose

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cores` | `integer` | Specifies vCores, which represent the logical CPU. Valid values are 2, 4, 8, 16, or 24, 32, 48, 80 | N | 2 |
| `storage` | `integer` | Specifies the amount of storage to allocate in GB. Ranges from 5 to 1048 | N | 5 |

Additional Provision Parameters for: business-critical

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cores` | `integer` | Specifies vCores, which represent the logical CPU. Valid values are 2, 4, 8, 16, or 24, 32, 48, 80 | N | 2 |
| `storage` | `integer` | Specifies the amount of storage to allocate in GB. Ranges from 5 to 1048 | N | 5 |

##### Bind

Creates a new user on the primary SQL Database. (The secondary database syncs the creation.) The new user will be named randomly and granted permission to log into and administer the database.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and credentials:

| Field Name | Type | Description |
|------------|------|-------------|
| `host` | `string` | The fully-qualified address of the Failover Group. |
| `port` | `int` | The port number to connect to on the SQL Server. |
| `database` | `string` | The name of the database. |
| `username` | `string` | The name of the database user. |
| `password` | `string` | The password for the database user. |
| `uri` | `string` | A uri string containing connection information. |
| `jdbcUrl` | `string` | A fully formed JDBC url. |
| `encrypt` | `boolean` | Flag indicating if the connection should be encrypted. |
| `tags` | `string[]` | List of tags. |

##### Update

Updates both the primary database and the secondary database.

###### Updating Parameters

Parameters for updating the SQL DB Database differ by plan. See each section for relevant parameters.

Additional Provision Parameters for : standard plan

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `dtus` | `integer` | Specifies Database transaction units, which represent a bundled measure of compute, storage, and IO resources. Valid values are 10, 20, 50, 100, 200, 400, 800, 1600, 3000 | N | 10 |

Additional Provision Parameters for : premium plan

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `dtus` | `integer` | Specifies Database transaction units, which represent a bundled measure of compute, storage, and IO resources. Valid values are 125, 250, 500, 1000, 1750, 4000 | N | 125 |

Additional Provision Parameters for: general-purpose

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cores` | `integer` | Specifies vCores, which represent the logical CPU. Valid values are 2, 4, 8, 16, or 24, 32, 48, 80 | N | 2 |
| `storage` | `integer` | Specifies the amount of storage to allocate in GB. Ranges from 5 to 1048. Note, decreasing storage is not currently supported | N | 5 |

Additional Provision Parameters for: business-critical

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cores` | `integer` | Specifies vCores, which represent the logical CPU. Valid values are 2, 4, 8, 16, or 24, 32, 48, 80 | N | 2 |
| `storage` | `integer` | Specifies the amount of storage to allocate in GB. Ranges from 5 to 1048. Note, decreasing storage is not currently supported | N | 5 |

##### Unbind

Drops the applicable user from the SQL Database.

##### Deprovision

Deletes the databases and the failover group.

##### Examples

###### Kubernetes

To add.

###### Cloud Foundry

Using the `cf` cli, you can create the `basic` plan of the `azure-sql-12-0-dr-database-pair` service with the following command:

```console
cf create-service azure-sql-database basic azure-sql-database -c '{
  "parentAlias":"serverpair",
  "database":"testdb",
  "failoverGroup":"testfg",
  "dtus": 100
}'
```

Note: this uses the alias used when provisioning the DBMS-only service above.

###### cURL

To provision an instance using the broker directly, you must use the ID for both plan and service. Assuming your OSBA is running locally on port 8080 with the default username and password, you can provision the `database` plan onto a previously provisioned Azure SQL Database DBMS with a cURL command similar to the following example. Note, this uses the alias provided in the DBMS-only example above:

```console
curl -X PUT \
  'http://localhost:8080/v2/service_instances/sql-database?accepts_incomplete=true' \
  -H 'authorization: Basic dXNlcm5hbWU6cGFzc3dvcmQ=' \
  -H 'content-type: application/json' \
  -H 'x-broker-api-version: 2.13' \
  -d '
{
  "service_id" : "2bbc160c-e279-4757-a6b6-4c0a4822d0aa",
  "plan_id" : "624828a9-c73c-4d35-bc9d-ea41cfc75853",
  "parameters" : {
    "parentAlias":"serverpair",
    "database":"testdb",
    "failoverGroup":"testfg",
    "dtus": 100
  }
}
'
```

### Service: azure-sql-12-0-dr-database-pair-registered

It is the service to *register* existing [**azure-sql-12-0-dr-database-pair** service](#service-azure-sql-12-0-dr-database-pair) instance as a service instance. It is for your OSBA instances in other regions to use the same failover group. It doesn't create new databases and doesn't delete databases but only validates the databases. Provisioning parameters can be referred to **azure-sql-12-0-dr-database-pair** service. Update is NOT supported.

### Service: azure-sql-12-0-dr-database-pair-from-existing

It is the service to *take over* existing [**azure-sql-12-0-dr-database-pair** service](#service-azure-sql-12-0-dr-database-pair) instance as a service instance. It is for migrating existing failover groups into OSBA's management. It doesn't create new databases in provisioning but deletes databases in deprovisioning. Provisioning parameters can be referred to **azure-sql-12-0-dr-database-pair** service. Update is also supported.

### Service: azure-sql-12-0-dr-database-pair-from-existing-primary

It is the service to create the secondary database and the failover group based on an existing primary database. It is for taking over an existing database and extending it to a failover group deployment. It creates new secondary database in provisioning and deletes both databases in deprovisioning. Provisioning parameters can be referred to **azure-sql-12-0-dr-database-pair** service. Update is also supported.
