# [Azure Text Analytics](https://azure.microsoft.com/en-us/services/cognitive-services/text-analytics/)

_Note: This module is EXPERIMENTAL and future releases may break the API._

## Services & Plans

### Service: azure-text-analytics

| Plan Name | Description |
|-----------|-------------|
| `free` | Free with 5,000 monthly transactions and no overage. |
| `standard-s0` | 25,000 monthly transactions and 3.00 per 1,000 overage. |
| `standard-s1` | 100,000 monthly transactions and 2.50 per 1,000 overage. |
| `standard-s2` | 500,000 monthly transactions and 2.00 per 1,000 overage. |
| `standard-s3` | 2,500,000 monthly transactions and 1.00 per 1,000 overage. |
| `standard-s4` | 10,000,000 monthly transactions and .50 per 1,000 overage. |

#### Behaviors

##### Provision

Provisions a new text analytics API.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Y |  |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | Y |  |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |

##### Bind

Returns the API address and access key.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details:

| Field Name | Type | Description |
|------------|------|-------------|
| `textAnalyticsEndpoint` | `string` | The text analytics API endpoint address. |
| `textAnalyticsKey` | `string` | The text analytics API access key. |
| `textAnalyticsName` | `string` | The name of the text analytics API. |

##### Unbind

Does nothing.

##### Deprovision

Deletes the text analytics API.