# Azure Container Instances

[Azure Container Instances](https://azure.microsoft.com/en-us/services/container-instances/) is a service for running Docker containers in Azure without needing to pay for or maintain VM resources. Just pay per second of container uptime.

## Warning!

This module's value is somewhat dubious, most platforms (e.g. Kubernetes or Cloud Foundry) that might integrate with this broker already have native mechanisms for running containers. (And in the case of Kubernetes, the ACI Connector even offers the option to seamlessly schedule containers in ACI.) Additionally, this module currently does not expose several critical options, which in all practicality are required to make this module useful. This would include, for instance, the ability to set environment variables within a container.__

__Do NOT depend on the service provided by this module in its current state.__

## Services & Plans

### azure-aci

| Plan Name | Description |
|-----------|-------------|
| `aci` | Run a container in ACI |

#### Behaviors

##### Provision
  
Provisions a new container in ACI.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Y | |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | N | A new resource group will be created with a UUID as its name. |
| `image` | `string` || Y ||
| `cpuCores` | `int` || N | `1` |
| `memoryInGb` | `float` || N | `1.5` |
| `port` | `int` | The port to expose. Currently exactly one must be exposed. It is an oversight that exposing neither 0 nor more than one is supported. This will be corrected soon. | Y ||
| `tags` | `object` | Tags to be applied to new resources, specified as key/value pairs. | N | |
  
##### Bind
  
Returns the public IP of the container.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details:

| Field Name | Type | Description |
|------------|------|-------------|
| `containerIPv4Address` | `string` | The container's public IP address. |

##### Unbind

Does nothing.
  
##### Deprovision

Deletes the container.
