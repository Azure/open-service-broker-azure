# [Azure Cognitive Services](https://azure.microsoft.com/en-us/services/cognitive-services/)

_Note: This module is EXPERIMENTAL and future releases may break the API._

## Services & Plans

### Service: azure-cognitive-services

| Plan Name | Description |
|-----------|-------------|
| `text-analytics` | Run a text analytics API with Azure Cognitive Services. |

#### Behaviors

##### Provision

Provisions a new text analytics API.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Required _unless_ an administrator has configured the broker itself with a default location. | The broker's default location, if configured. |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | N | If an administrator has configured the broker itself with a default resource group and node is specified, that default will be applied, otherwise, a new resource group will be created with a UUID as its name. |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |

##### Bind

Returns the API address and access key.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details:

| Field Name | Type | Description |
|------------|------|-------------|
| `Endpoint` | `string` | The text analytics API endpoint address. |
| `CognitiveKey` | `string` | The text analytics API access key. |

##### Unbind

Does nothing.

##### Deprovision

Deletes the text analytics API.
