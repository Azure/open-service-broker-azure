# [Azure Database for PostgreSQL](https://azure.microsoft.com/en-us/services/postgresql/)

|![](https://upload.wikimedia.org/wikipedia/commons/thumb/1/17/Warning.svg/50px-Warning.svg.png) | This module is EXPERIMENTAL. It is under heavy development and remains subject to the possibility of breaking changes. |
|---|---|

## Services & Plans

### Service: azure-postgresqldb

| Plan Name | Description |
|-----------|-------------|
| `basic50` | Basic Tier, 50 DTUs |
| `basic100` | Basic Tier, 100 DTUs |

#### Behaviors

##### Provision
  
Provisions a new PostgreSQL server and a new database upon that server. The new
database will be named randomly and will be owned by a role (group) of the same
name.

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
  
Creates a new role (user) on the PostgreSQL server. The new role will be named
randomly and added to the  role (group) that owns the database.

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
