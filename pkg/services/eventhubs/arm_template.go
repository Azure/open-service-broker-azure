package eventhubs

// nolint: lll
var armTemplateBytes = []byte(`
{
	"$schema": "http://schema.management.azure.com/schemas/2014-04-01-preview/deploymentTemplate.json#",
	"contentVersion": "1.0.0.0",
	"parameters": {
		"tags": {
			"type": "object"
		}
	},
	"variables": {
		"defaultSASKeyName": "RootManageSharedAccessKey",
		"authRuleResourceId": "[resourceId('Microsoft.EventHub/namespaces/authorizationRules', '{{.eventHubNamespace}}', variables('defaultSASKeyName'))]",
		"ehVersion": "2017-04-01"
	},
	"resources": [
		{
			"apiVersion": "2017-04-01",
			"name": "{{.eventHubNamespace}}",
			"type": "Microsoft.EventHub/Namespaces",
			"location": "{{.location}}",
			"sku": {
				"name": "{{.eventHubSku}}"
			},
			"tags": "[parameters('tags')]",
			"resources": [
				{
					"apiVersion": "2017-04-01",
					"name": "{{.eventHubName}}",
					"type": "EventHubs",
					"dependsOn": [
						"Microsoft.EventHub/namespaces/{{.eventHubNamespace}}"
					],
					"properties": {
						"messageRetentionInDays": 1,
						"partitionCount": 4
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
