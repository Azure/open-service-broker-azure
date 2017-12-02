# Module Name Here

_Describe the Azure service(s) provided by the module, in general terms. This is a good place for links to relevant Azure documentation._

## Services & Plans

### Service Name Here

_Describe the service. Repeat for each if the module provides multiple services._

_Describe the plans of the service using a table like this:_

| Plan Name | Description |
|-----------|-------------|
| `Plan Name Here` | _Description_ |

#### Behaviors

_Describe the behaviors of the provision, bind, unbind, and deprovision operations, as implemented by this module (for this service)._

##### Provision
  
_Describe what is provisioned. (i.e. What Azure resources are created?)_

###### Provisioning Parameters

_Describe the supported parameters using a table like this:_

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Y | |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | N | A new resource group will be created with a UUID as its name. |
| `tags` | `object` | Tags to be applied to new resources, specified as key/value pairs. | N | |
| `Parameter Name Here` | `Type` | _Description_ | Y/N |  `Default value` |
  
##### Bind
  
_Describe what occurs upon binding. (i.e. Is a new user account of some kind
created? Or is there a single set of credentials that are shared upon bind?)_

###### Binding Parameters

_Describe the supported parameters using a table like this:_

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `Parameter Name Here` | `Type` | _Description_ | Y/N |  `Default value` |


###### Credentials

_Describe the schema of any credentials that are returned by the binding
operation using a table like this:_

| Field Name | Type | Description |
|------------|------|-------------|
| `Field Name Here` | `Type` | _Description_ |

##### Unbind

_Describe what occurs upon unbinding. (i.e. Is a user account deleted?)_
  
##### Deprovision

_Describe what occurs upon deprovisioning. (i.e. What Azure resources are deleted?)_
