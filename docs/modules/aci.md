# [Azure Container Instances](https://azure.microsoft.com/en-us/services/container-instances/)

_Note: This module is EXPERIMENTAL and is not included in the General Availability release of Open Service Broker for Azure. It will be added in a future OSBA release._

## Services & Plans

### Service: azure-aci

| Plan Name | Description |
|-----------|-------------|
| `aci` | Run a container in ACI |

#### Behaviors

##### Provision
  
Provisions a new container in ACI.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `cpuCores` | `int` | The number of virtual CPU cores requested for the container. | N | `1` |
| `image` | `string` | The Docker image on which to base the container. | Y ||
| `location` | `string` | The Azure region in which to provision applicable resources. | Required _unless_ an administrator has configured the broker itself with a default location. | The broker's default location, if configured. |
| `memoryInGb` | `float64` | Gigabytes of memory requested for the container. Must be specified in increments of 0.10 GB. | N | `1.5` |
| `ports` | `[]int` | The port(s) to open on the container. The container will be assigned a public IP (v4) address if and only if one or more ports are opened. | Y ||
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | N | If an administrator has configured the broker itself with a default resource group and nonde is specified, that default will be applied, otherwise, a new resource group will be created with a UUID as its name. |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |
  
##### Bind
  
Returns the public IP of the container _only_ if the container exposes one or
more ports.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details:

| Field Name | Type | Description |
|------------|------|-------------|
| `publicIPv4Address` | `string` | The container's public IP (v4) address. Note that this field is returned upon bind _only_ if the container exposes one or more ports. |

##### Unbind

Does nothing.
  
##### Deprovision

Deletes the container.
