# [Azure Database for PostgreSQL](https://azure.microsoft.com/en-us/services/postgresql/)

|![](https://upload.wikimedia.org/wikipedia/commons/thumb/1/17/Warning.svg/50px-Warning.svg.png) | This module is EXPERIMENTAL. It is under heavy development and remains subject to the possibility of breaking changes. |
|---|---|

Open Service Broker for Azure contains three Azure Database for PostgreSQL services. These services enable you to select the most appropriate provision scenario for your needs. These services are:

| Service Name | Description |
|--------------|-------------|
| `azure-postgresql` | Provision both an Azure Database for PostgreSQL Database Management System (DBMS) and a database. |
| `azure-postgresql-dbms-only` | Provision only an Azure Database for PostgreSQL DBMS. This can be used to provision multiple databases at a later time. |
| `azure-postgresql-database-only` | Provision a new database only upon a previously provisioned DBMS. |

The `azure-postgresql` service allows you to provision both a DBMS and a database. This service is ready to use upon successful provisioning. You can not provision additional databases onto an instance provisioned through this service. The `azure-postgresql-dbms-only` and `azure-postgresql-database-only` services, on the other hand, can be combined to provision multiple databases on a single DBMS.  For more information on each service, refer to the descriptions below.

## Services & Plans

### Service: azure-postgresql

| Plan Name | Description |
|-----------|-------------|
| `basic50` | Basic Tier, 50 DTUs |
| `basic100` | Basic Tier, 100 DTUs |

#### Behaviors

##### Provision

Provisions a new PostgreSQL DBMS and a new database upon that DBMS. The new
database will be named randomly and will be owned by a role (group) of the same
name.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Required _unless_ an administrator has configured the broker itself with a default location. | The broker's default location, if configured. |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | N | If an administrator has configured the broker itself with a default resource group and none is specified, that default will be applied, otherwise, a new resource group will be created with a UUID as its name. |
| `sslEnforcement` | `string` | Specifies whether the server requires the use of TLS when connecting. Valid valued are `""` (unspecified), `enabled`, or `disabled`. | N | `""`. Left unspecified, SSL _will_ be enforced. |
| `firewallRules`  | `array` | Specifies the firewall rules to apply to the server. Definition follows. | N | `[]` Left unspecified, Firewall will default to only Azure IPs. If rules are provided, they must have valid values. |
| `firewallRules[n].name` | `string` | Specifies the name of the generated firewall rule |Y | |
| `firewallRules[n].startIPAddress` | `string` | Specifies the start of the IP range allowed by this firewall rule | Y | |
| `firewallRules[n].endIPAddress` | `string` | Specifies the end of the IP range allowed by this firewall rule | Y | |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |
| `extensions` | `string[]` | Specifies a list of PostgreSQL extensions to install | N | |

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
| `sslRequired` | `boolean` | Flag indicating if SSL is required to connect the MySQL DBMS. |
| `uri` | `string` | A URI string containing all necessary connection information. |
| `tags` | `string[]` | A list of tags consumers can use to identify the credential. |

##### Unbind

Drops the applicable role (user) from the PostgreSQL DBMS.

##### Deprovision

Deletes the PostgreSQL DBMS and database.

##### Examples

###### Kubernetes

The `contrib/k8s/examples/postgresql/postgresql-all-in-one-instance.yaml` can be used to provision the `basic50` plan. This can be done with the following example:

```console
kubectl create -f contrib/k8s/examples/postgresql/postgresql-all-in-one-instance.yaml
```

You can then create a binding with the following command:

```console
kubectl create -f contrib/k8s/examples/postgresql/postgresql-all-in-one-binding.yaml
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

To provision an instance using the broker directly, you must use the ID for both plan and service. Assuming your OSBA is running locally on port 8080 with the default username and password, you can provision the `basic50` plan with a cUrl command similar to the following example:

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

### Service: azure-postgresql-dbms-only

| Plan Name | Description |
|-----------|-------------|
| `basic50` | Basic Tier, 50 DTUs |
| `basic100` | Basic Tier, 100 DTUs |

#### Behaviors

##### Provision

Provisions an Azure Database for PostgreSQL DBMS instance containing no databases. Databases can be created through subsequent provision requests using the `azure-postgresql-database-only` service.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Required _unless_ an administrator has configured the broker itself with a default location. | The broker's default location, if configured. |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | N | If an administrator has configured the broker itself with a default resource group and none is specified, that default will be applied, otherwise, a new resource group will be created with a UUID as its name. |
| `alias` | `string` | Specifies an alias that can be used by later provision actions to create databases on this DBMS. | Y | |
| `sslEnforcement` | `string` | Specifies whether the server requires the use of TLS when connecting. Valid valued are `""` (unspecified), `enabled`, or `disabled`. | N | `""`. Left unspecified, SSL _will_ be enforced. |
| `firewallRules`  | `array` | Specifies the firewall rules to apply to the server. Definition follows. | N | `[]` Left unspecified, Firewall will default to only Azure IPs. If rules are provided, they must have valid values. |
| `firewallRules[n].name` | `string` | Specifies the name of the generated firewall rule |Y | |
| `firewallRules[n].startIPAddress` | `string` | Specifies the start of the IP range allowed by this firewall rule | Y | |
| `firewallRules[n].endIPAddress` | `string` | Specifies the end of the IP range allowed by this firewall rule | Y | |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |

##### Bind

This service is not bindable.

##### Unbind

This service is not bindable.

##### Deprovision

Deletes the PostgreSQL DBMS only. If databases have been provisioned on this DBMS, deprovisioning will be deferred until all databases have been deprovisioned.

##### Examples

###### Kubernetes

The `contrib/k8s/examples/postgresql/postgresql-dbms-only-instance.yaml` can be used to provision the `basic50` plan. This can be done with the following example:

```console
kubectl create -f contrib/k8s/examples/postgresql/postgresql-dbms-only-instance.yaml
```

###### Cloud Foundry

Using the `cf` cli, you can provision the `basic50` plan of this service with the following command:

```console
cf create-service azure-postgresql-dbms-only basic50 postgresql-dbms-only -c '{
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

To provision an instance using the broker directly, you must use the ID for both plan and service. Assuming your OSBA is running locally on port 8080 with the default username and password, you can provision the `basic50` plan with a cUrl command similar to the following example:

```console
curl -X PUT \
  'http://localhost:8080/v2/service_instances/postgreqsl-dbms-only?accepts_incomplete=true' \
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

### Service: azure-postgresql-database-only

| Plan Name | Description |
|-----------|-------------|
| `database-only` | New database on existing DBMS |

#### Behaviors

##### Provision

Provisions a new PostgreSQL DBMS and a new database upon that DBMS. The new
database will be named randomly and will be owned by a role (group) of the same
name.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `extensions` | `string[]` | Specifies a list of PostgreSQL extensions to install | N | |
| `parentAlias` | `string` | Specifies the alias of the DBMS upon which the database should be provisioned. | Y | |

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
| `sslRequired` | `boolean` | Flag indicating if SSL is required to connect the MySQL DBMS. |
| `uri` | `string` | A URI string containing all necessary connection information. |
| `tags` | `string[]` | A list of tags consumers can use to identify the credential. |

##### Unbind

Drops the applicable role (user) from the PostgreSQL DBMS.

##### Deprovision

Deletes the PostgreSQL database only, the DBMS remains provisioned.

##### Examples

###### Kubernetes

The `contrib/k8s/examples/postgresql/postgresql-database-only-instance.yaml` can be used to provision the `database-only` plan. This can be done with the following example:

```console
kubectl create -f contrib/k8s/examples/postgresql/postgresql-database-only-instance.yaml
```

You can then create a binding with the following command:

```console
kubectl create -f contrib/k8s/examples/postgresql/postgresql-database-only-binding.yaml
```

###### Cloud Foundry

Using the `cf` cli, you can provision the `database-only` plan of this service with the following command:

```console
cf create-service azure-postgresql-database-only database-only postgresql-db-only -c '{
    "parentAlias" : "ed9798f2-2e91-4b21-8903-d364a3ff7d12"
}'
```

###### cURL

To provision an instance using the broker directly, you must use the ID for both plan and service. Assuming your OSBA is running locally on port 8080 with the default username and password, you can provision the `database-only` plan with a cUrl command similar to the following example:

```console
curl -X PUT \
  'http://localhost:8080/v2/service_instances/postgresql-db-only?accepts_incomplete=true' \
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