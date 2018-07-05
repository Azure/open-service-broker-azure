package servicebus

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
		"defaultAuthRuleResourceId": "[resourceId('Microsoft.ServiceBus/namespaces/authorizationRules', '{{.serviceBusNamespaceName}}', variables('defaultSASKeyName'))]",
		"sbVersion": "2017-04-01"
	},
	"resources": [
		{
			"apiVersion": "2017-04-01",
			"name": "{{.serviceBusNamespaceName}}",
			"type": "Microsoft.ServiceBus/Namespaces",
			"location": "{{.location}}",
			"sku": {
				"name": "{{.serviceBusSku}}"
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
