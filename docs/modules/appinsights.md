# [Azure Application Insights](https://docs.microsoft.com/en-us/azure/application-insights/app-insights-overview)

_Note: This module is EXPERIMENTAL. To enable this module, you must run Open Service Broker for Azure with modules.minStability set to `experimental`_

## Services & Plans

### Service: azure-appinsights

| Plan Name | Description |
|-----------|-------------|
| `default` | For general application. |

#### Behaviors

##### Provision

Provisions a new Application Insights.

###### Provisioning Parameters

| Parameter Name | Type | Description | Required | Default Value |
|----------------|------|-------------|----------|---------------|
| `location` | `string` | The Azure region in which to provision applicable resources. | Y |  |
| `resourceGroup` | `string` | The (new or existing) resource group with which to associate new resources. | Y |  |
| `appInsightsName` | `string` | The Application Insights component name. | N | A randomly generated UUID. |
| `tags` | `map[string]string` | Tags to be applied to new resources, specified as key/value pairs. | N | Tags (even if none are specified) are automatically supplemented with `heritage: open-service-broker-azure`. |

##### Bind

Returns the instrumentation key.

###### Binding Parameters

This binding operation does not support any parameters.

###### Credentials

Binding returns the following connection details and shared credentials:

| Field Name | Type | Description |
|------------|------|-------------|
| `instrumentationKey` | `string` | Instrumentation key. |

##### Unbind

Does nothing.

##### Deprovision

Deletes the Application Insights.
