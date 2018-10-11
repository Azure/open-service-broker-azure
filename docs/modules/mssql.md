# [Azure SQL Database](https://azure.microsoft.com/en-us/services/sql-database/)

Open Service Broker for Azure (OSBA) contains a number of Azure SQL Database services. These services enable you to select the most appropriate provisioning scenario for your needs. These services are:

| Service Name | Description |
|--------------|-------------|
| [`azure-sql-12-0`](#service-azure-sql-12-0) | Provision both a SQL Server DBMS and a database. |
| [`azure-sql-12-0-dbms`](#service-azure-sql-12-0-dbms) | Provision only a SQL Server Database Management System (DBMS). This can be used to provision multiple databases at a later time. |
| [`azure-sql-12-0-database`](#service-azure-sql-12-0-database) | Provision a new database only upon a previously provisioned DBMS. |
| [`azure-sql-12-0-dbms-registered`](#service-azure-sql-12-0-dbms-registered) | Register an existing server as a DBMS service instance. |
| [`azure-sql-12-0-database-from-existing`](#service-azure-sql-12-0-database-from-existing) | Taking over an existing database upon a previous DBMS as a database service instance. The service requires `ENABLE_MIGRATION_SERVICES` to be `true` in OSBA environment variables. |

The `azure-sql-12-0` service allows you to provision both a DBMS and a database. When the provision operation is successful, the database will be ready to use. You can not provision additional databases onto an instance provisioned through this service. The `azure-sql-12-0-dbms` and `azure-sql-12-0-database` services, on the other hand, can be combined to provision multiple databases on a single DBMS. The `azure-sql-12-0-dbms-registered` is in the same position with `azure-sql-12-0-dbms`. And `azure-sql-12-0-database-from-existing` is in the same position with `azure-sql-12-0-database`. For more information on each service, refer to the descriptions below.

## Services & Plans

### Service: azure-sql-12-0

| Plan Name | Description |
|-----------|-------------|
| `basic` | Basic Tier, 5 DTUs, 2GB, 7 days point-in-time restore |
| `standard` | Standard Tier, Up to 3000 DTUs, with 250GB storage, 35 days point-in-time restore |
| `premium` | Premium Tier, Up to 4000 DTUs, with 500GB storage, 35 days point-in-time restore |
| `general-purpose` | General Purpose Tier, Up to 80 vCores, Up to 440 GB Memory, Up to 1 TB storage, 7 days point-in-time restore |
| `business-critical` | Business Critical Tier, Up to 80 vCores, Up to 440 GB Memory, Up to 1 TB storage, Local SSD, 7 days point-in-time restore. Offers highest resilience to failures using several isolated replicas |

For applications which require less than 1 core, please use the `basic`, `standard` or `premium` plans.

#### Behaviors

##### Provision

Provisions a new SQL Server and a new database upon that server. The new dbms and database will be named randomly.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Y | |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | Y | |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |
| `firewallRules`  | `array` | Specifies the firewall rules to apply to the server. Definition follows. | N | `[]` Left unspecified, Firewall will default to only Azure IPs. If rules are provided, they must have valid values. |
| `firewallRules[n].name` | `string` | Specifies the name of the generated firewall rule |Y | |
| `firewallRules[n].startIPAddress` | `string` | Specifies the start of the IP range allowed by this firewall rule | Y | |
| `firewallRules[n].endIPAddress` | `string` | Specifies the end of the IP range allowed by this firewall rule | Y | |

Additional Provision Parameters for : standard plan

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `dtus` | `integer` | Specifies Database transaction units, which represent a bundled measure of compute, storage, and IO resources. Valid values are 10, 20, 50, 100, 200, 400, 800, 1600, 3000 | N | 10 |

Additional Provision Parameters for : premium plan

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `dtus` | `integer` | Specifies Database transaction units, which represent a bundled measure of compute, storage, and IO resources. Valid values are 125, 250, 500, 1000, 1750, 1000 | N | 125 |

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
| `uri` | `string` | A uri string containing connection information. |
| `jdbcUrl` | `string` | A fully formed JDBC url. |
| `encrypt` | `boolean` | Flag indicating if the connection should be encrypted. |
| `tags` | `string[]` | List of tags. |

##### Update

Updates a previously provisioned SQL DB Database and DBMS.

###### Updating Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `firewallRules`  | `array` | Specifies the firewall rules to apply to the server. Definition follows. | N | `[]` Left unspecified, Firewall will default to only Azure IPs. If rules are provided, they must have valid values. |
| `firewallRules[n].name` | `string` | Specifies the name of the generated firewall rule |Y | |
| `firewallRules[n].startIPAddress` | `string` | Specifies the start of the IP range allowed by this firewall rule | Y | |
| `firewallRules[n].endIPAddress` | `string` | Specifies the end of the IP range allowed by this firewall rule | Y | |

Parameters for updating the SQL DB Database differ by plan. See each section for relevant parameters.

Additional Provision Parameters for : standard plan

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `dtu` | `integer` | Specifies Database transaction units, which represent a bundled measure of compute, storage, and IO resources. Valid values are 10, 20, 50, 100, 200, 400, 800, 1600, 3000 | N | 10 |

Additional Provision Parameters for : premium plan

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `dtu` | `integer` | Specifies Database transaction units, which represent a bundled measure of compute, storage, and IO resources. Valid values are 125, 250, 500, 1000, 1750, 1000 | N | 125 |

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

Drops the applicable user from the SQL Server.

##### Deprovision

Deletes both the database and the SQL Server instance.

##### Examples

###### Kubernetes

The `contrib/k8s/examples/sql/sql-instance.yaml` can be used to provision one of the plans from the all-in-one `azure-sql` service. This can be done with the following example:

```console
kubectl create -f ../../contrib/k8s/examples/sql/sql-instance.yaml
```

You can then create a binding to the service with the following command:

```console
kubectl create -f ../../contrib/k8s/examples/sql/sql-binding.yaml
```

###### Cloud Foundry

Using the `cf` cli, you can create the `basic` plan of the `azure-sql` service with the following command:

```console
cf create-service azure-sql-12-0 basic azure-sql-all-in-one -c '{
        "resourceGroup" : "demo",
        "location" : "eastus",
        "firewallRules" : [
            {
                "name": "AllowAll",
                "startIPAddress": "0.0.0.0",
                "endIPAddress" : "255.255.255.255"
            }
        ]
    }
'
```

###### cURL

Assuming your OSBA is running locally on port 8080 with the default username and password, you can provision the all-in-one Azure SQL Database service with a cURL command similar to the following example:

```console
curl -X PUT \
  'http://localhost:8080/v2/service_instances/azure-sql-12-0?accepts_incomplete=true' \
  -H 'authorization: Basic dXNlcm5hbWU6cGFzc3dvcmQ=' \
  -H 'content-type: application/json' \
  -H 'x-broker-api-version: 2.13' \
  -d '{
    "service_id" : "fb9bc99e-0aa9-11e6-8a8a-000d3a002ed5",
    "plan_id" : "17725188-76a2-4d6c-8e86-49f146766eeb",
    "parameters" : {
        "resourceGroup" : "demo",
        "location" : "eastus",
        "firewallRules" : [
            {
                "name": "AllowAll",
                "startIPAddress": "0.0.0.0",
                "endIPAddress" : "255.255.255.255"
            }
        ]
    }
}'
```

### Service: azure-sql-12-0-dbms

| Plan Name | Description |
|-----------|-------------|
| `dbms` | Azure SQL Server DBMS-Only |

#### Behaviors

##### Provision

Provisions a SQL Server DBMS instance containing no databases. Databases can be created through subsequent provision requests using the `azure-sql-12-0-database` service.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Y | |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | Y | |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |
| `alias` | `string` | Specifies an alias that can be used by later provision actions to create databases on this DBMS. | Y | |
| `firewallRules`  | `array` | Specifies the firewall rules to apply to the server. Definition follows. | N | `[]` Left unspecified, Firewall will default to only Azure IPs. If rules are provided, they must have valid values. |
| `firewallRules[n].name` | `string` | Specifies the name of the generated firewall rule |Y | |
| `firewallRules[n].startIPAddress` | `string` | Specifies the start of the IP range allowed by this firewall rule | Y | |
| `firewallRules[n].endIPAddress` | `string` | Specifies the end of the IP range allowed by this firewall rule | Y | |

##### Bind

This service is not bindable.

##### Update

Updates a previously provisioned SQL DB DBMS.

###### Updating Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `firewallRules`  | `array` | Specifies the firewall rules to apply to the server. Definition follows. | N | `[]` Left unspecified, Firewall will default to only Azure IPs. If rules are provided, they must have valid values. |
| `firewallRules[n].name` | `string` | Specifies the name of the generated firewall rule |Y | |
| `firewallRules[n].startIPAddress` | `string` | Specifies the start of the IP range allowed by this firewall rule | Y | |
| `firewallRules[n].endIPAddress` | `string` | Specifies the end of the IP range allowed by this firewall rule | Y | |

##### Unbind

This service is not bindable.

##### Deprovision

Deprovision will delete the SQL Server DBMS. If any databases have been provisioned on this DBMS, deprovisioning will be deferred until all databases have been deprovisioned.

##### Examples

###### Kubernetes

The `contrib/k8s/examples/sql/advanced/sql-dbms-instance.yaml` can be used to provision one of the plans from the `azure-sql-12-0-dbms` service. This can be done with the following example:

```console
kubectl create -f ../../contrib/k8s/examples/sql/advanced/sql-dbms-instance.yaml
```

###### Cloud Foundry

Using the `cf` cli, you can create the `dbms` plan of the `azure-sql-12-0-dbms` service with the following command:

```console
cf create-service azure-sql-12-0-dbms dbms azure-sql-dbms -c '{
        "resourceGroup" : "demo",
        "location" : "eastus",
        "alias" : "ed9798f2-2e91-4b21-8903-d364a3ff7d12",
        "firewallRules" : [
            {
                "name": "AllowAll",
                "startIPAddress": "0.0.0.0",
                "endIPAddress" : "255.255.255.255"
            }
        ]
    }
'
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
        "resourceGroup": "demo",
        "location" : "eastus",
        "alias" : "7ce1fc7d-073a-4be6-abe6-536e7053496d",
        "firewallRules" : [
            {
                "name": "AllowSome",
                "startIPAddress": "0.0.0.0",
                "endIPAddress" : "35.0.0.0"
            },
            {
                "name": "AllowMore",
                "startIPAddress": "35.0.0.1",
                "endIPAddress" : "255.255.255.255"
            }
        ]
    }
}'
```

### Service: azure-sql-12-0-database

| Plan Name | Description |
|-----------|-------------|
| `basic` | Basic Tier, 5 DTUs, 2GB, 7 days point-in-time restore |
| `standard` | Standard Tier, Up to 3000 DTUs, with 250GB storage, 35 days point-in-time restore |
| `premium` | Premium Tier, Up to 4000 DTUs, with 500GB storage, 35 days point-in-time restore |
| `general-purpose` | General Purpose Tier, Up to 80 vCores, Up to 440 GB Memory, Up to 1 TB storage, 7 days point-in-time restore |
| `business-critical` | Business Critical Tier, Up to 80 vCores, Up to 440 GB Memory, Up to 1 TB storage, Local SSD, 7 days point-in-time restore. Offers highest resilience to failures using several isolated replicas |

#### Behaviors

##### Provision

Provisions a new database upon an existing server. The new database will be named randomly. If the DBMS does not yet exist, provision of the database will be deferred until the DBMS has been provisioned.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `parentAlias` | `string` | Specifies the alias of the DBMS upon which the database should be provisioned. | Y | |

Additional Provision Parameters for : standard plan

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `dtus` | `integer` | Specifies Database transaction units, which represent a bundled measure of compute, storage, and IO resources. Valid values are 10, 20, 50, 100, 200, 400, 800, 1600, 3000 | N | 10 |


Additional Provision Parameters for : premium plan

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `dtus` | `integer` | Specifies Database transaction units, which represent a bundled measure of compute, storage, and IO resources. Valid values are 125, 250, 500, 1000, 1750, 1000 | N | 125 |

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
| `uri` | `string` | A uri string containing connection information. |
| `jdbcUrl` | `string` | A fully formed JDBC url. |
| `encrypt` | `boolean` | Flag indicating if the connection should be encrypted. |
| `tags` | `string[]` | List of tags. |

##### Update

Updates a previously provisioned SQL DB Database.

###### Updating Parameters

Parameters for updating the SQL DB Database differ by plan. See each section for relevant parameters.

Additional Provision Parameters for : standard plan

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `dtu` | `integer` | Specifies Database transaction units, which represent a bundled measure of compute, storage, and IO resources. Valid values are 10, 20, 50, 100, 200, 400, 800, 1600, 3000 | N | 10 |

Additional Provision Parameters for : premium plan

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `dtu` | `integer` | Specifies Database transaction units, which represent a bundled measure of compute, storage, and IO resources. Valid values are 125, 250, 500, 1000, 1750, 1000 | N | 125 |

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

Drops the applicable user from the SQL Server.

##### Deprovision

Deletes the database from the SQL Server instance, but does not delete the DBMS.

##### Examples

###### Kubernetes

The `contrib/k8s/examples/sql/advanced/sql-database-instance.yaml` can be used to provision the `basic` plan. This can be done with the following example:

```console
kubectl create -f ../../contrib/k8s/examples/sql/advanced/sql-database-instance.yaml
```
You can then create a binding with the following command:

```console
kubectl create -f ../../contrib/k8s/examples/sql/advanced/sql-database-binding.yaml
```

###### Cloud Foundry

Using the `cf` cli, you can create the `basic` plan of the `azure-sql-12-0-database` service with the following command:

```console
cf create-service azure-sql-12-0-database basic azure-sql-database -c '{
    "parentAlias" : "ed9798f2-2e91-4b21-8903-d364a3ff7d12"
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
        "parentAlias" : "7ce1fc7d-073a-4be6-abe6-536e7053496d"
    }
}
'
```

### Service: azure-sql-12-0-dbms-registered

It is to *register* an existing Azure SQL Server as a SQL DBMS service instance. Both **azure-sql-12-0-database** service and **azure-sql-12-0-database-from-existing** service can be its child service.

#### Behaviors

##### Provision

Please refer to [**azure-sql-12-0-dbms** service](#service-azure-sql-12-0-dbms) with extra required provisioning parameters below.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `server` | `string` | The SQL server name. | Y | |
| `administratorLogin` | `string` | The administratorLogin input when creating the SQL server. | Y | |
| `administratorLoginPassword` | `string` | The administratorLoginPassword input when creating the SQL server. | Y | |

##### Update

Update the `administratorLogin` and/or `administratorLoginPassword` as they may change and the server is assumed to be managed by yourself.

###### Updating Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `administratorLogin` | `string` | New administratorLogin. | N | |
| `administratorLoginPassword` | `string` | New administratorLoginPassword. | N | |

##### Bind, Unbind

Please refer to [**azure-sql-12-0-dbms** service](#service-azure-sql-12-0-dbms).

##### Deprovision

Do nothing. The SQL server would not be deleted in Azure.

### Service: azure-sql-12-0-database-from-existing

It is to create SQL database service instance from existing Azure SQL Database *for taking over the database*. Typically, you can create **azure-sql-12-0-dbms-registered** service instance to register your Azure SQL server first and use this service to import the database to OSBA's management.

##### Provision

Please refer to [**azure-sql-12-0-database** service](#service-azure-sql-12-0-database) with extra required provisioning parameters below.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `database` | `string` | The SQL database name. | Y | |

##### Update, Bind, Unbind, Deprovision

Please refer to [**azure-sql-12-0-database** service](#service-azure-sql-12-0-database).

Note that deprovision would delete the SQL database in Azure!
