package cosmosdb

// nolint: lll
var armTemplateBytes = []byte(`
{
	"$schema": "http://schema.management.azure.com/schemas/2014-04-01-preview/deploymentTemplate.json#",
	"contentVersion": "1.0.0.0",
	"parameters": {
		"location": {
			"type": "string"
		},
		"name": {
			"type": "string"
		},
		"kind": {
			"type": "string"
		},
		"tags": {
			"type": "object"
		}
	},
	"resources": [
		{
			"apiVersion": "2015-04-08",
			"kind": "[parameters('kind')]",
			"type": "Microsoft.DocumentDb/databaseAccounts",
			"name": "[parameters('name')]",
			"location": "[parameters('location')]",
			"properties": {
				"databaseAccountOfferType": "Standard",
				"locations": [
					{
						"id": "[concat(parameters('name'), '-', parameters('location'))]",
						"failoverPriority": 0,
						"locationName": "[parameters('location')]"
					}
				]
			},
			"tags": "[parameters('tags')]"
		}
	],
	"outputs": {
		"fullyQualifiedDomainName": {
			"type": "string",
			"value": "[reference(parameters('name')).documentEndpoint]"
		},
		"primaryKey":{
			"type": "string",
			"value": "[listKeys(resourceId('Microsoft.DocumentDb/databaseAccounts', parameters('name')), '2015-04-08').primaryMasterKey]"
		}
	}
}
`)
