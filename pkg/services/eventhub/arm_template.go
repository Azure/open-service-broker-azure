package eventhub

// nolint: lll
var armTemplateBytes = []byte(`
{
	"$schema": "http://schema.management.azure.com/schemas/2014-04-01-preview/deploymentTemplate.json#",
	"contentVersion": "1.0.0.0",
	"parameters": {
		"eventHubNamespace": {
			"type": "string",
			"metadata": {
				"description": "Name of the EventHub namespace"
			}
		},
		"eventHubName": {
			"type": "string",
			"metadata": {
				"description": "Name of the Event Hub"
			}
		},
		"messageRetentionInDays": {
			"type": "int",
			"defaultValue": 1,
			"minValue": 1,
			"maxValue": 7,
			"metadata": {
				"description": "How long to retain the data in Event Hub"
			}
		},
		"partitionCount": {
			"type": "int",
			"defaultValue": 4,
			"minValue": 2,
			"maxValue": 32,
			"metadata": {
				"description": "Number of partitions chosen"
			}
		},
		"eventHubSku": {
			"type": "string",
			"allowedValues": [
				"Basic",
				"Standard"
			],
			"metadata": {
				"description": "Tiers for Event Hubs"
			}
		},
		"tags": {
			"type": "object"
		}
	},
	"variables": {
		"defaultSASKeyName": "RootManageSharedAccessKey",
		"authRuleResourceId": "[resourceId('Microsoft.EventHub/namespaces/authorizationRules', parameters('eventHubNamespace'), variables('defaultSASKeyName'))]",
		"ehVersion": "2017-04-01"
	},
	"resources": [
		{
			"apiVersion": "2017-04-01",
			"name": "[parameters('eventHubNamespace')]",
			"type": "Microsoft.EventHub/Namespaces",
			"location": "[resourceGroup().location]",
			"sku": {
				"name": "[parameters('eventHubSku')]"
			},
			"tags": "[parameters('tags')]",
			"resources": [
				{
					"apiVersion": "2017-04-01",
					"name": "[parameters('eventHubName')]",
					"type": "EventHubs",
					"dependsOn": [
						"[concat('Microsoft.EventHub/namespaces/', parameters('eventHubNamespace'))]"
					],
					"properties": {
						"messageRetentionInDays": "[parameters('messageRetentionInDays')]",
						"partitionCount": "[parameters('partitionCount')]"
					},
					"tags": "[parameters('tags')]"
				}
			]
		}
	],
	"outputs": {
		"ConnectionString": {
			"type": "string",
			"value": "[listkeys(variables('authRuleResourceId'), variables('ehVersion')).primaryConnectionString]"
		},
		"PrimaryKey": {
			"type": "string",
			"value": "[listkeys(variables('authRuleResourceId'), variables('ehVersion')).primaryKey]"
		}
	}
}
`)
