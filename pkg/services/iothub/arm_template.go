package iothub

// nolint: lll
var armTemplateBytes = []byte(`
{
    "$schema": "http://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json",
    "contentVersion": "1.0.0.0",
    "parameters": {
        "tags": {
            "type": "object"
        }
    },
    "resources": [
        {
            "type": "Microsoft.Devices/IotHubs",
            "sku": {
                "name": "{{.skuName}}",
                "capacity": "{{.skuUnits}}"
            },
            "name": "{{.iotHubName}}",
            "apiVersion": "2018-04-01",
            "location": "{{.location}}",
            "properties": {
                "eventHubEndpoints": {
                    "events": {
                        "retentionTimeInDays": 1,
                        "partitionCount": "{{.partitionCount}}"
                    }
                }
            }
        }
    ],
    "outputs": {
        "keyInfo": {
            "type": "Object",
            "value": "[listKeys(resourceId('Microsoft.Devices/iotHubs', '{{.iotHubName}}'), '2018-04-01').value[0]]"
        }
    }
}
`)
