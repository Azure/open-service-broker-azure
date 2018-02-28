# Module Name Here

_Describe the Azure service(s) provided by the module, in general terms. This is
a good place for links to relevant Azure documentation._

_Include or exclude the following disclaimer, as appropriate. For brand new
modules, it should be included._

|![](https://upload.wikimedia.org/wikipedia/commons/thumb/1/17/Warning.svg/50px-Warning.svg.png) | This module is EXPERIMENTAL. It is under heavy development and remains subject to the possibility of breaking changes. |
|---|---|

## Services & Plans

### Service: Service Name Here

_Describe the service. Repeat for each if the module provides multiple
services._

_Describe the plans of the service using a table like this:_

| Plan Name | Description |
|-----------|-------------|
| `Plan Name Here` | _Description_ |

#### Behaviors

_Describe the behaviors of the provision, bind, unbind, and deprovision
operations, as implemented by this module (for this service)._

##### Provision
  
_Describe what is provisioned. (i.e. What Azure resources are created?)_

###### Provisioning Parameters

_Describe the supported parameters using a table like this:_

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Required _unless_ an administrator has configured the broker itself with a default location. | The broker's default location, if configured. |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | N | If an administrator has configured the broker itself with a default resource group and nonde is specified, that default will be applied, otherwise, a new resource group will be created with a UUID as its name. |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |
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

_Describe what occurs upon deprovisioning. (i.e. What Azure resources are
deleted?)_

##### Examples

###### Kubernetes

_Provide an example with manifests. Manifests should be added to contrib/k8s/examples/{module}_

###### Cloud Foundry

_Provide an example with using cf cli_

###### cURL

_Provide an example of using the broker directly_