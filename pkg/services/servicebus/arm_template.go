package servicebus

// nolint: lll
var armTemplateBytes = []byte(`
{
	"$schema": "http://schema.management.azure.com/schemas/2014-04-01-preview/deploymentTemplate.json#",
	"contentVersion": "1.0.0.0",
	"parameters": {
		"serviceBusNamespaceName": {
			"type": "string",
			"metadata": {
				"description": "Name of the Service Bus namespace"
			}
		},
		"serviceBusSku": {
			"type": "string",
			"allowedValues": [
				"Basic",
				"Standard",
				"Premium"
			],
			"metadata": {
				"description": "The messaging tier for service Bus namespace"
			}
		},
		"tags": {
			"type": "object"
		}
	},
	"variables": {
		"defaultSASKeyName": "RootManageSharedAccessKey",
		"defaultAuthRuleResourceId": "[resourceId('Microsoft.ServiceBus/namespaces/authorizationRules', parameters('serviceBusNamespaceName'), variables('defaultSASKeyName'))]",
		"sbVersion": "2017-04-01"
	},
	"resources": [
		{
			"apiVersion": "2017-04-01",
			"name": "[parameters('serviceBusNamespaceName')]",
			"type": "Microsoft.ServiceBus/Namespaces",
			"location": "[resourceGroup().location]",
			"sku": {
				"name": "[parameters('serviceBusSku')]"
			},
			"tags": "[parameters('tags')]"
		}
	],
	"outputs": {
		"ConnectionString": {
			"type": "string",
			"value": "[listkeys(variables('defaultAuthRuleResourceId'), variables('sbVersion')).primaryConnectionString]"
		},
		"PrimaryKey": {
			"type": "string",
			"value": "[listkeys(variables('defaultAuthRuleResourceId'), variables('sbVersion')).primaryKey]"
		}
	}
}
`)
