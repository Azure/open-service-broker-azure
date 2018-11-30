# [Azure Database for MySQL](https://azure.microsoft.com/en-us/services/mysql/)

Open Service Broker for Azure contains three Azure Database for MySQL services. These services enable you to select the most appropriate provision scenario for your needs. These services are:

| Service Name | Description |
|--------------|-------------|
| `azure-mysql-5-7` | Provision both an Azure Database for MySQL Database Management System (DBMS) and a database, using MySQL 5.7 |
| `azure-mysql-5-7-dbms` | Provision only an Azure Database for MySQL DBMS with MySQL 5.7. This can be used to provision multiple databases at a later time. |
| `azure-mysql-5-7-database` | Provision a new database only upon a previously provisioned DBMS. |

The `azure-mysql-5-7` service allows you to provision both a DBMS and a database. When the provision operation is successful, the database will be ready to use. You can't provision additional databases onto an instance provisioned through this service. The `azure-mysql-5-7-dbms` and `azure-mysql-5-7-database` services, on the other hand, can be combined to provision multiple databases on a single DBMS.  For more information on each service, refer to the descriptions below.

_This module involves the Parent-Child Model concept in OSBA, please refer to the [Parent-Child Model doc](../parent-child-model-for-multiple-layers-services.md)_.

## Services & Plans

### Service: azure-mysql-5-7

| Plan Name | Description |
|-----------|-------------|
| `basic` | Basic Tier, Up to 2 vCores, Variable I/O performance |
| `general-purpose` | General Purporse Tier, Up to 32 vCores, Predictable I/O Performance, Local or Geo-Redundant Backups |
| `memory-optimized` | Memory Optimized Tier, Up to 16 memory optimized vCores, Predictable I/O Performance, Local or Geo-Redundant Backups |

#### Behaviors

##### Provision

Provisions a new MySQL DBMS and a new database upon it. The new database will be named randomly.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Y | |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | Y | |
| `sslEnforcement` | `string` | Specifies whether the server requires the use of TLS when connecting. Valid valued are `""` (unspecified), `enabled`, or `disabled`. | N | `""`. Left unspecified, SSL _will_ be enforced. |
| `firewallRules`  | `array` | Specifies the firewall rules to apply to the server. Definition follows. | N | `[]` Left unspecified, Firewall will default to only Azure IPs. If rules are provided, they must have valid values. |
| `firewallRules[n].name` | `string` | Specifies the name of the generated firewall rule |Y | |
| `firewallRules[n].startIPAddress` | `string` | Specifies the start of the IP range allowed by this firewall rule | Y | |
| `firewallRules[n].endIPAddress` | `string` | Specifies the end of the IP range allowed by this firewall rule | Y | |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |

The three plans each have additional provisioning parameters with different default and allowed values. See the tables below for details on each.

Provisioning Parameters: basic

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cores` | `integer` | Specifies vCores, which represent the logical CPU. Valid values are 1 or 2 | N | 1 |
| `storage` | `integer` | Specifies the amount of storage to allocate in GB. Ranges from 5 to 1048 | N | 10 |
| `backupRetention` | `integer` | Specifies the number of days to retain backups. Ranges from 7 to 35 | N | 7 |

Provisioning Parameters: general-purpose

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cores` | `integer` | Specifies vCores, which represent the logical CPU. Valid values are 2, 4, 8, 16 or 32 | N | 2 |
| `storage` | `integer` | Specifies the amount of storage to allocate in GB. Ranges from 5 to 2048 | N | 10 |
| `backupRetention` | `integer` | Specifies the number of days to retain backups. Ranges from 7 to 35 | N | 7 |
| `backupRedundancy` | `string` | Specifies the backup redundancy, either `local` or `geo` | N | `local` |

Provisioning Parameters: memory-optimized

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cores` | `integer` | Specifies vCores, which represent the logical CPU. Valid values are 2, 4, 8 or 16 | N | 2 |
| `storage` | `integer` | Specifies the amount of storage to allocate in GB. Ranges from 5 to 2048 | N | 10 |
| `backupRetention` | `integer` | Specifies the number of days to retain backups. Ranges from 7 to 35 | N | 7 |
| `backupRedundancy` | `string` | Specifies the backup redundancy, either `local` or `geo` | N | `local` |

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
| `sslRequired` | `boolean` | Flag indicating if SSL is required to connect the MySQL DBMS. |
| `uri` | `string` | A URI string containing all necessary connection information. |
| `tags` | `string[]` | A list of tags consumers can use to identify the credential. |

##### Unbind

Drops the applicable user from the MySQL DBMS.

##### Deprovision

Deletes the MySQL DBMS and database.

##### Examples

###### Kubernetes

The `contrib/k8s/examples/mysql/mysql-instance.yaml` can be used to provision the `basic` plan. This can be done with the following example:

```console
kubectl create -f contrib/k8s/examples/mysql/mysql-instance.yaml
```

You can then create a binding with the following command:

```console
kubectl create -f contrib/k8s/examples/mysql/mysql-binding.yaml
```

###### Cloud Foundry

Using the `cf` cli, you can provision the `basic` plan of this service with the following command:

```console
cf create-service azure-mysql basic mysql-all-in-one -c '{
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

To provision an instance using the broker directly, you must use the ID for both plan and service. Assuming your OSBA is running locally on port 8080 with the default username and password, you can provision the `basic` plan with a cURL command similar to the following example:

```console
curl -X PUT \
  'http://localhost:8080/v2/service_instances/mysql-all-in-one?accepts_incomplete=true' \
  -H 'authorization: Basic dXNlcm5hbWU6cGFzc3dvcmQ=' \
  -H 'content-type: application/json' \
  -H 'x-broker-api-version: 2.13' \
  -d '{
    "service_id" : "997b8372-8dac-40ac-ae65-758b4a5075a5",
    "plan_id" : "427559f1-bf2a-45d3-8844-32374a3e58aa",
    "parameters" : {
        "resourceGroup": "demo",
        "location" : "eastus",
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

### Service: azure-mysql-5-7-dbms

| Plan Name | Description |
|-----------|-------------|
| `basic` | Basic Tier, Up to 2 vCores, Variable I/O performance |
| `general-purpose` | General Purporse Tier, Up to 32 vCores, Predictable I/O Performance, Local or Geo-Redundant Backups |
| `memory-optimized` | Memory Optimized Tier, Up to 16 memory optimized vCores, Predictable I/O Performance, Local or Geo-Redundant Backups |

#### Behaviors

##### Provision

Provisions an Azure Database for MySQL DBMS instance containing no databases. Databases can be created through subsequent provision requests using the `azure-mysql-database` service.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Y | |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | Y | |
| `alias` | `string` | Specifies an alias that can be used by later provision actions to create databases on this DBMS. | Y | |
| `sslEnforcement` | `string` | Specifies whether the server requires the use of TLS when connecting. Valid valued are `""` (unspecified), `enabled`, or `disabled`. | N | `""`. Left unspecified, SSL _will_ be enforced. |
| `firewallRules`  | `array` | Specifies the firewall rules to apply to the server. Definition follows. | N | `[]` Left unspecified, Firewall will default to only Azure IPs. If rules are provided, they must have valid values. |
| `firewallRules[n].name` | `string` | Specifies the name of the generated firewall rule |Y | |
| `firewallRules[n].startIPAddress` | `string` | Specifies the start of the IP range allowed by this firewall rule | Y | |
| `firewallRules[n].endIPAddress` | `string` | Specifies the end of the IP range allowed by this firewall rule | Y | |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |

The three plans each have additional provisioning parameters with different default and allowed values. See the tables below for details on each.

Provisioning Parameters: basic

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cores` | `integer` | Specifies vCores, which represent the logical CPU. Valid values are 1 or 2 | N | 1 |
| `storage` | `integer` | Specifies the amount of storage to allocate in GB. Ranges from 5 to 1048 | N | 10 |
| `backupRetention` | `integer` | Specifies the number of days to retain backups. Ranges from 7 to 35 | N | 7 |

Provisioning Parameters: general-purpose

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cores` | `integer` | Specifies vCores, which represent the logical CPU. Valid values are 2, 4, 8, 16 or 32 | N | 2 |
| `storage` | `integer` | Specifies the amount of storage to allocate in GB. Ranges from 5 to 2048 | N | 10 |
| `backupRetention` | `integer` | Specifies the number of days to retain backups. Ranges from 7 to 35 | N | 7 |
| `backupRedundancy` | `string` | Specifies the backup redundancy, either `local` or `geo` | N | `local` |

Provisioning Parameters: memory-optimized

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cores` | `integer` | Specifies vCores, which represent the logical CPU. Valid values are 2, 4, 8 or 16 | N | 2 |
| `storage` | `integer` | Specifies the amount of storage to allocate in GB. Ranges from 5 to 2048 | N | 10 |
| `backupRetention` | `integer` | Specifies the number of days to retain backups. Ranges from 7 to 35 | N | 7 |
| `backupRedundancy` | `string` | Specifies the backup redundancy, either `local` or `geo` | N | `local` |

##### Bind

This service is not bindable.

##### Unbind

This service is not bindable.

##### Deprovision

Deprovision will delete the MySQL DBMS. If any databases have been provisioned on this DBMS, deprovisioning will be deferred until all databases have been deprovisioned.

##### Examples

###### Kubernetes

The `contrib/k8s/examples/mysql/advanced/mysql-dbms-instance.yaml` can be used to provision the `basic50` plan. This can be done with the following example:

```console
kubectl create -f contrib/k8s/examples/mysql/advanced/mysql-dbms-instance.yaml
```

###### Cloud Foundry

Using the `cf` cli, you can provision the `basic50` plan of this service with the following command:

```console
cf create-service azure-mysql-dbms basic50 mysql-dbms -c '{
    "resourceGroup" : "demo",
    "location" : "eastus",
    "alias" : "679aab6d-39e7-4a41-8b45-49975569079c",
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

To provision an instance using the broker directly, you must use the ID for both plan and service. Assuming your OSBA is running locally on port 8080 with the default username and password, you can provision the `basic50` plan with a cURL command similar to the following example:

```console
curl -X PUT \
  'http://localhost:8080/v2/service_instances/mysql-dbms?accepts_incomplete=true' \
  -H 'authorization: Basic dXNlcm5hbWU6cGFzc3dvcmQ=' \
  -H 'content-type: application/json' \
  -H 'x-broker-api-version: 2.13' \
  -d '{
    "service_id" : "30e7b836-199d-4335-b83d-adc7d23a95c2",
    "plan_id" : "3f65ebf9-ac1d-4e77-b9bf-918889a4482b",
    "parameters" : {
        "resourceGroup": "demo",
        "location" : "eastus",
        "alias" : "c33c3b3a-c491-4197-8ced-66d4f89baa67",
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

### Service: azure-mysql-5-7-dbms

| Plan Name | Description |
|-----------|-------------|
| `database` | New database on existing MySQL DBMS |

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
| `sslRequired` | `boolean` | Flag indicating if SSL is required to connect the MySQL DBMS. |
| `uri` | `string` | A URI string containing all necessary connection information. |
| `tags` | `string[]` | A list of tags consumers can use to identify the credential. |

##### Unbind

Drops the applicable user from the MySQL DBMS.

##### Deprovision

Deletes the database from the MySQL DBMS. The DBMS itself is not deprovisioned.

##### Examples

###### Kubernetes

The `contrib/k8s/examples/mysql/advanced/mysql-database-instance.yaml` can be used to provision the `database` plan. This can be done with the following example:

```console
kubectl create -f contrib/k8s/examples/mysql/advanced/mysql-database-instance.yaml
```

You can then create a binding with the following command:

```console
kubectl create -f contrib/k8s/examples/mysql/advanced/mysql-database-binding.yaml
```

###### Cloud Foundry

Using the `cf` cli, you can provision the `database` plan of this service with the following command:

```console
cf create-service azure-mysql-databasey database mysql-database -c '{
    "parentAlias" : "679aab6d-39e7-4a41-8b45-49975569079c"
}
'
```

Note: this uses the alias provided in the DBMS-only example.

###### cURL

To provision an instance using the broker directly, you must use the ID for both plan and service. Assuming your OSBA is running locally on port 8080 with the default username and password, you can provision the `database` plan with a cURL command similar to the following example. Note, this uses the alias provided in the DBMS-only example above:

```console
curl -X PUT \
  'http://localhost:8080/v2/service_instances/mysql-database?accepts_incomplete=true' \
  -H 'authorization: Basic dXNlcm5hbWU6cGFzc3dvcmQ=' \
  -H 'content-type: application/json' \
  -H 'x-broker-api-version: 2.13' \
  -d '{
    "service_id" : "6704ae59-3eae-49e9-82b4-4cbcc00edf08",
    "plan_id" : "ec77bd04-2107-408e-8fde-8100c1ce1f46",
    "parameters" : {
        "parentAlias" : "c33c3b3a-c491-4197-8ced-66d4f89baa67"
    }
}'
```
