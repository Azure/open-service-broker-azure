# [Azure Database for PostgreSQL](https://azure.microsoft.com/en-us/services/postgresql/)

*Note: PostgreSQL version 9.6 is STABLE, and PostgreSQL version 10 is in PREVIEW*

Open Service Broker for Azure contains three types of Azure Database for PostgreSQL services. These services enable you to select the most appropriate provision scenario for your needs. These services are:

| Service Type                  | Description                                                  |
| ----------------------------- | ------------------------------------------------------------ |
| `azure-postgresql-*`          | Provision both an Azure Database for PostgreSQL Database Management System (DBMS) and a database. |
| `azure-postgresql-*-dbms`     | Provision only an Azure Database for PostgreSQL DBMS. This can be used to provision multiple databases at a later time. |
| `azure-postgresql-*-database` | Provision a new database only upon a previously provisioned DBMS. |

The `azure-postgresql-*` services allow you to provision both a DBMS and a database. When the provision operation is successful, the database will be ready to use. You can not provision additional databases onto an instance of `azure-postgresql-*`. The `azure-postgresql-*-dbms` and `azure-postgresql-*-database` services, on the other hand, can be combined to provision multiple databases on a single DBMS. Currently, OSBA supports two versions of Azure Database for PostgreSQL services:
<table>
	<thead>
		<tr>
			<th>Service Type</th>
			<th>Service name</th>
			<th>Stability</th>
		</tr>
	</thead>
	<tbody>
		<tr>
			<td rowspan=2>azure-postgresql-*</td>
			<td>azure-postgresql-9-6</td>
			<td>Stable</td>
		</tr>
		<tr>
			<td>azure-postgresql-10</td>
			<td>Preview</td>
		</tr>
		<tr>
			<td rowspan=2>azure-postgresql-*-dbms</td>
			<td>azure-postgresql-9-6-dbms</td>
			<td>Stable</td>
		</tr>
		<tr>
			<td>azure-postgresql-10-dbms</td>
			<td>Preview</td>
		</tr>
		<tr>
			<td rowspan=2>azure-postgresql-*-database</td>
			<td>azure-postgresql-9-6-database</td>
			<td>Stable</td>
		</tr>
		<tr>
			<td>azure-postgresql-10-database</td>
			<td>Preview</td>
		</tr>
	</tbody>
</table>

 For more information on each service, refer to the descriptions below.

_This module involves the Parent-Child Model concept in OSBA, please refer to the [Parent-Child Model doc](../parent-child-model-for-multiple-layers-services.md)_.

## Services & Plans

### Service: azure-postgresql-*

| Plan Name | Description |
|-----------|-------------|
| `basic` | Basic Tier, Up to 2 vCores, Variable I/O performance |
| `general-purpose` | General Purporse Tier, Up to 32 vCores, Predictable I/O Performance, Local or Geo-Redundant Backups |
| `memory-optimized` | Memory Optimized Tier, Up to 16 memory optimized vCores, Predictable I/O Performance, Local or Geo-Redundant Backups |

#### Behaviors

##### Provision

Provisions a new PostgreSQL DBMS and a new database upon that DBMS. The new
database will be named randomly and will be owned by a role (group) of the same
name.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Y | |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | Y | |
| `serverName` | `string` | Name of the PostgreSQL server. | N | A random generated string. |
| `adminAccountSettings` | `object` | Settings of administrator account of PostgreSQL server. Typically you do not need to specify this. | N | Default admin username is "postgres" and password is a randomly generated string. |
| `adminAccountSettings.adminUsername` | `string` | The administrator username for the server. | N | "postgres" |
| `adminAccountSettings.adminPassword` | `string` | The administrator password for the server. **Warning**: you may leak your password if you specify this property, others can see this password in your request body and `ServiceInstance` definition. DO NOT use this property unless you know what you are doing. | N | A random generated password. |
| `sslEnforcement` | `string` | Specifies whether the server requires the use of TLS when connecting. Valid valued are `""` (unspecified), `enabled`, or `disabled`. | N | `""`. Left unspecified, SSL _will_ be enforced. |
| `firewallRules`  | `array` | Specifies the firewall rules to apply to the server. Definition follows. | N | `[]` Left unspecified, Firewall will default to only Azure IPs. If rules are provided, they must have valid values. |
| `firewallRules[n].name` | `string` | Specifies the name of the generated firewall rule |Y | |
| `firewallRules[n].startIPAddress` | `string` | Specifies the start of the IP range allowed by this firewall rule | Y | |
| `firewallRules[n].endIPAddress` | `string` | Specifies the end of the IP range allowed by this firewall rule | Y | |
| `virtualNetworkRules`  | `array` | Specifies the firewall rules to apply to the server. Definition follows. | N | `[]` Left unspecified, Firewall will default to only Azure IPs. If rules are provided, they must have valid values. |
| `virtualNetworkRules[n].name` | `string` | Specifies the name of the generated virtual network rule |Y | |
| `virtualNetworkRules[n].subnetId` | `string` | The full resource ID of a subnet in a virtual network to allow access from. Example format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/virtualNetworks/{vn}/subnets/{sn} | Y | |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |
| `extensions` | `string[]` | Specifies a list of PostgreSQL extensions to install | N | |

####### Provisioning Parameters: basic

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cores` | `integer` | Specifies vCores, which represent the logical CPU. Valid values are 1 or 2 | N | 1 |
| `storage` | `integer` | Specifies the amount of storage to allocate in GB. Ranges from 5 to 1048 | N | 10 |
| `backupRetention` | `integer` | Specifies the number of days to retain backups. Ranges from 7 to 35 | N | 7 |


####### Provisioning Parameters: general-purpose

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cores` | `integer` | Specifies vCores, which represent the logical CPU. Valid values are 2, 4, 8, 16, or 32 | N | 2 |
| `storage` | `integer` | Specifies the amount of storage to allocate in GB. Ranges from 5 to 1048 | N | 10 |
| `backupRetention` | `integer` | Specifies the number of days to retain backups. Ranges from 7 to 35 | N | 7 |
| `backupRedundancy` | `string` | Specifies the backup redundancy, either `local` or `geo` | N | `local` |

####### Provisioning Parameters: memory-optimized

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cores` | `integer` | Specifies vCores, which represent the logical CPU. Valid values are 2, 4, 8, 16 | N | 2 |
| `storage` | `integer` | Specifies the amount of storage to allocate in GB. Ranges from 5 to 1048 | N | 10 |
| `backupRetention` | `integer` | Specifies the number of days to retain backups. Ranges from 7 to 35 | N | 7 |
| `backupRedundancy` | `string` | Specifies the backup redundancy, either `local` or `geo` | N | `local` |

##### Update

Updates a previously provisioned PostgreSQL DBMS. Currently updating the database extensions is not supported.

###### Updating Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `sslEnforcement` | `string` | Specifies whether the server requires the use of TLS when connecting. Valid valued are `""` (unspecified), `enabled`, or `disabled`. | N | `""`. Left unspecified, SSL _will_ be enforced. |
| `firewallRules`  | `array` | Specifies the firewall rules to apply to the server. Definition follows. | N | `[]` Left unspecified, Firewall will default to only Azure IPs. If rules are provided, they must have valid values. |
| `firewallRules[n].name` | `string` | Specifies the name of the generated firewall rule |Y | |
| `firewallRules[n].startIPAddress` | `string` | Specifies the start of the IP range allowed by this firewall rule | Y | |
| `firewallRules[n].endIPAddress` | `string` | Specifies the end of the IP range allowed by this firewall rule | Y | |

####### Updating Parameters: basic

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cores` | `integer` | Specifies vCores, which represent the logical CPU. Valid values are 1 or 2 | N | 1 |
| `storage` | `integer` | Specifies the amount of storage to allocate in GB. Ranges from 5 to 1048. Note, this must not be lower than what was given at provision time. | N | 10 |
| `backupRetention` | `integer` | Specifies the number of days to retain backups. Ranges from 7 to 35 | N | 7 |


####### Updating Parameters: general-purpose

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cores` | `integer` | Specifies vCores, which represent the logical CPU. Valid values are 2, 4, 8, 16, or 32 | N | 2 |
| `storage` | `integer` | Specifies the amount of storage to allocate in GB. Ranges from 5 to 1048. Note, this must not be lower than what was given at provision time. | N | 10 |
| `backupRetention` | `integer` | Specifies the number of days to retain backups. Ranges from 7 to 35 | N | 7 |

####### Updating Parameters: memory-optimized

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cores` | `integer` | Specifies vCores, which represent the logical CPU. Valid values are 2, 4, 8, or 16 | N | 2 |
| `storage` | `integer` | Specifies the amount of storage to allocate in GB. Ranges from 5 to 1048. Note, this must not be lower than what was given at provision time. | N | 10 |
| `backupRetention` | `integer` | Specifies the number of days to retain backups. Ranges from 7 to 35 | N | 7 |

##### Bind

Creates a new role (user) on the PostgreSQL DBMS. The new role will be named
randomly and added to the  role (group) that owns the database.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and credentials:

| Field Name | Type | Description |
|------------|------|-------------|
| `host` | `string` | The fully-qualified address of the PostgreSQL DBMS. |
| `port` | `int` | The port number to connect to on the PostgreSQL DBMS. |
| `database` | `string` | The name of the database. |
| `username` | `string` | The name of the database user (in the form username@host). |
| `password` | `string` | The password for the database user. |
| `sslRequired` | `boolean` | Flag indicating if SSL is required to connect the PostgreSQL DBMS. |
| `uri` | `string` | A URI string containing all necessary connection information. |
| `jdbcUrl` | `string` | A fully formed JDBC url. |
| `tags` | `string[]` | A list of tags consumers can use to identify the credential. |

##### Unbind

Drops the applicable role (user) from the PostgreSQL DBMS.

##### Deprovision

Deletes the PostgreSQL DBMS and database.

##### Examples

###### Kubernetes

The `contrib/k8s/examples/postgresql/postgresql-instance.yaml` can be used to provision the `basic50` plan. This can be done with the following example:

```console
kubectl create -f contrib/k8s/examples/postgresql/postgresql-instance.yaml
```

You can then create a binding with the following command:

```console
kubectl create -f contrib/k8s/examples/postgresql/postgresql-binding.yaml
```

###### Cloud Foundry

Using the `cf` cli, you can provision the `basic50` plan of this service with the following command:

```console
cf create-service azure-postgresql basic50 postgresql-all-in-one -c '{
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

To provision an instance using the broker directly, you must use the ID for both plan and service. Assuming your OSBA is running locally on port 8080 with the default username and password, you can provision the `basic50` plan with a cURL command similar to the following example:

```console
curl -X PUT \
  'http://localhost:8080/v2/service_instances/postgresql-all-in-one?accepts_incomplete=true' \
  -H 'authorization: Basic dXNlcm5hbWU6cGFzc3dvcmQ=' \
  -H 'content-type: application/json' \
  -H 'x-broker-api-version: 2.13' \
  -d '{
    "service_id" : "b43b4bba-5741-4d98-a10b-17dc5cee0175",
    "plan_id" : "b2ed210f-6a10-4593-a6c4-964e6b6fad62",
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

### Service: azure-postgresql-*-dbms

| Plan Name | Description |
|-----------|-------------|
| `basic` | Basic Tier, Up to 2 vCores, Variable I/O performance |
| `general-purpose` | General Purporse Tier, Up to 32 vCores, Predictable I/O Performance, Local or Geo-Redundant Backups |
| `memory-optimized` | Memory Optimized Tier, Up to 16 memory optimized vCores, Predictable I/O Performance, Local or Geo-Redundant Backups |

#### Behaviors

##### Provision

Provisions an Azure Database for PostgreSQL DBMS instance containing no databases. Databases can be created through subsequent provision requests using the `azure-postgresql-database` service.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Y | |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | Y | |
| `serverName` | `string` | Name of the PostgreSQL server. | N | A random generated string. |
| `adminUsername` | `string` | The administrator username for the server. | N | "postgresql" |
| `adminPassword` | `string` | The administrator password for the server. **Warning**: you may leak your password if you specify this property, others can see this password in your request body and `ServiceInstance` definition. DO NOT use this property unless you know what you are doing. | N | A random generated password. |
| `alias` | `string` | Specifies an alias that can be used by later provision actions to create databases on this DBMS. | Y | |
| `sslEnforcement` | `string` | Specifies whether the server requires the use of TLS when connecting. Valid valued are `""` (unspecified), `enabled`, or `disabled`. | N | `""`. Left unspecified, SSL _will_ be enforced. |
| `firewallRules`  | `array` | Specifies the firewall rules to apply to the server. Definition follows. | N | `[]` Left unspecified, Firewall will default to only Azure IPs. If rules are provided, they must have valid values. |
| `firewallRules[n].name` | `string` | Specifies the name of the generated firewall rule |Y | |
| `firewallRules[n].startIPAddress` | `string` | Specifies the start of the IP range allowed by this firewall rule | Y | |
| `firewallRules[n].endIPAddress` | `string` | Specifies the end of the IP range allowed by this firewall rule | Y | |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |

####### Provisioning Parameters: basic

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cores` | `integer` | Specifies vCores, which represent the logical CPU. Valid values are 1 or 2 | N | 1 |
| `storage` | `integer` | Specifies the amount of storage to allocate in GB. Ranges from 5 to 1048 | N | 10 |
| `backupRetention` | `integer` | Specifies the number of days to retain backups. Ranges from 7 to 35 | N | 7 |


####### Provisioning Parameters: general-purpose

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cores` | `integer` | Specifies vCores, which represent the logical CPU. Valid values are 2, 4, 8, 16, or 32 | N | 2 |
| `storage` | `integer` | Specifies the amount of storage to allocate in GB. Ranges from 5 to 1048 | N | 10 |
| `backupRetention` | `integer` | Specifies the number of days to retain backups. Ranges from 7 to 35 | N | 7 |
| `backupRedundancy` | `string` | Specifies the backup redundancy, either `local` or `geo` | N | `local` |

####### Provisioning Parameters: memory-optimized

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cores` | `integer` | Specifies vCores, which represent the logical CPU. Valid values are 2, 4, 8, or 16 | N | 2 |
| `storage` | `integer` | Specifies the amount of storage to allocate in GB. Ranges from 5 to 1048 | N | 10 |
| `backupRetention` | `integer` | Specifies the number of days to retain backups. Ranges from 7 to 35 | N | 7 |
| `backupRedundancy` | `string` | Specifies the backup redundancy, either `local` or `geo` | N | `local` |

##### Update

Updates a previously provisioned PostgreSQL DBMS.

###### Updating Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `sslEnforcement` | `string` | Specifies whether the server requires the use of TLS when connecting. Valid valued are `""` (unspecified), `enabled`, or `disabled`. | N | `""`. Left unspecified, SSL _will_ be enforced. |
| `firewallRules`  | `array` | Specifies the firewall rules to apply to the server. Definition follows. | N | `[]` Left unspecified, Firewall will default to only Azure IPs. If rules are provided, they must have valid values. |
| `firewallRules[n].name` | `string` | Specifies the name of the generated firewall rule |Y | |
| `firewallRules[n].startIPAddress` | `string` | Specifies the start of the IP range allowed by this firewall rule | Y | |
| `firewallRules[n].endIPAddress` | `string` | Specifies the end of the IP range allowed by this firewall rule | Y | |

####### Updating Parameters: basic

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cores` | `integer` | Specifies vCores, which represent the logical CPU. Valid values are 1 or 2 | N | 1 |
| `storage` | `integer` | Specifies the amount of storage to allocate in GB. Ranges from 5 to 1048. Note, this must not be lower than what was given at provision time. | N | 10 |
| `backupRetention` | `integer` | Specifies the number of days to retain backups. Ranges from 7 to 35 | N | 7 |


####### Updating Parameters: general-purpose

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cores` | `integer` | Specifies vCores, which represent the logical CPU. Valid values are 2, 4, 8, 16, or 32 | N | 2 |
| `storage` | `integer` | Specifies the amount of storage to allocate in GB. Ranges from 5 to 1048. Note, this must not be lower than what was given at provision time. | N | 10 |
| `backupRetention` | `integer` | Specifies the number of days to retain backups. Ranges from 7 to 35 | N | 7 |

####### Updating Parameters: memory-optimized

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cores` | `integer` | Specifies vCores, which represent the logical CPU. Valid values are 2, 4, 8, or 16 | N | 2 |
| `storage` | `integer` | Specifies the amount of storage to allocate in GB. Ranges from 5 to 1048. Note, this must not be lower than what was given at provision time. | N | 10 |
| `backupRetention` | `integer` | Specifies the number of days to retain backups. Ranges from 7 to 35 | N | 7 |

##### Bind

This service is not bindable.

##### Unbind

This service is not bindable.

##### Deprovision

Deletes the PostgreSQL DBMS only. If databases have been provisioned on this DBMS, deprovisioning will be deferred until all databases have been deprovisioned.

##### Examples

###### Kubernetes

The `contrib/k8s/examples/postgresql/advanced/postgresql-dbms-instance.yaml` can be used to provision the `basic50` plan. This can be done with the following example:

```console
kubectl create -f contrib/k8s/examples/postgresql/advanced/postgresql-dbms-instance.yaml
```

###### Cloud Foundry

Using the `cf` cli, you can provision the `basic50` plan of this service with the following command:

```console
cf create-service azure-postgresql-dbms basic50 postgresql-dbms -c '{
    "resourceGroup" : "demo",
    "location" : "eastus",
    "alias" : "3f368072-6fa8-42ad-ae9c-c02e59b7dc8d",
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
  'http://localhost:8080/v2/service_instances/postgreqsl-dbms?accepts_incomplete=true' \
  -H 'authorization: Basic dXNlcm5hbWU6cGFzc3dvcmQ=' \
  -H 'content-type: application/json' \
  -H 'x-broker-api-version: 2.13' \
  -d '{
    "service_id" : "d3f74b44-79bc-4d1e-bf7d-c247c2b851f9",
    "plan_id" : "bf389028-8dcc-433a-ab6f-0ee9b8db142f",
    "parameters" : {
        "resourceGroup": "demo",
        "location" : "eastus",
        "alias" : "d94f7740-74d8-446a-bbfd-c616935b4d58",
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

### Service: azure-postgresql-*-database

| Plan Name | Description |
|-----------|-------------|
| `database` | New database on existing DBMS |

#### Behaviors

##### Provision

Provisions a new database upon an existing PostgreSQL DBMS. The new
database will be named randomly and will be owned by a role (group) of the same
name.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `extensions` | `string[]` | Specifies a list of PostgreSQL extensions to install | N | |
| `parentAlias` | `string` | Specifies the alias of the DBMS upon which the database should be provisioned. | Y | |

**Note**: You should use corresponding  `dbms` service instance as the parent of `database` service instance.  For example, you should use `azure-postgresql-10-dbms` as the parent of `azure-postgresql-10-database`.

##### Update

Not currently supported.

##### Bind

Creates a new role (user) on the PostgreSQL DBNS. The new role will be named
randomly and added to the  role (group) that owns the database.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and credentials:

| Field Name | Type | Description |
|------------|------|-------------|
| `host` | `string` | The fully-qualified address of the PostgreSQL DBMS. |
| `port` | `int` | The port number to connect to on the PostgreSQL DBMS. |
| `database` | `string` | The name of the database. |
| `username` | `string` | The name of the database user (in the form username@host). |
| `password` | `string` | The password for the database user. |
| `sslRequired` | `boolean` | Flag indicating if SSL is required to connect the PostgreSQL DBMS. |
| `uri` | `string` | A URI string containing all necessary connection information. |
| `jdbcUrl` | `string` | A fully formed JDBC url. |
| `tags` | `string[]` | A list of tags consumers can use to identify the credential. |

##### Unbind

Drops the applicable role (user) from the PostgreSQL DBMS.

##### Deprovision

Deletes the PostgreSQL database only, the DBMS remains provisioned.

##### Examples

###### Kubernetes

The `contrib/k8s/examples/postgresql/postgresql-database-instance.yaml` can be used to provision the `database` plan. This can be done with the following example:

```console
kubectl create -f contrib/k8s/examples/postgresql/advanced/postgresql-database-instance.yaml
```

You can then create a binding with the following command:

```console
kubectl create -f contrib/k8s/examples/postgresql/advanced/postgresql-database-binding.yaml
```

###### Cloud Foundry

Using the `cf` cli, you can provision the `database` plan of this service with the following command:

```console
cf create-service azure-postgresql-database database postgresql-database -c '{
    "parentAlias" : "ed9798f2-2e91-4b21-8903-d364a3ff7d12"
}'
```

###### cURL

To provision an instance using the broker directly, you must use the ID for both plan and service. Assuming your OSBA is running locally on port 8080 with the default username and password, you can provision the `database` plan with a cURL command similar to the following example:

```console
curl -X PUT \
  'http://localhost:8080/v2/service_instances/postgresql-db?accepts_incomplete=true' \
  -H 'authorization: Basic dXNlcm5hbWU6cGFzc3dvcmQ=' \
  -H 'content-type: application/json' \
  -H 'x-broker-api-version: 2.13' \
  -d '{
    "service_id" : "25434f16-d762-41c7-bbdd-8045d7f74ca6",
    "plan_id" : "df6f5ef1-e602-406b-ba73-09c107d1e31b",
    "parameters" : {
        "parentAlias" : "d94f7740-74d8-446a-bbfd-c616935b4d58"
    }
}'
```
