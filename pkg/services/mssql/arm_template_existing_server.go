package mssql

// nolint: lll
var armTemplateExistingServerBytes = []byte(`
{
	"$schema": "http://schema.management.azure.com/schemas/2014-04-01-preview/deploymentTemplate.json#",
	"contentVersion": "1.0.0.0",
	"parameters": {
		"serverName": {
			"type": "string"
		},
		"databaseName": {
			"type": "string"
		},
		"edition": {
			"type": "string"
		},
		"requestedServiceObjectiveName": {
			"type": "string"
		},
		"maxSizeBytes": {
			"type": "string"
		},
		"tags": {
			"type": "object"
		}
	},
	"variables": {
		"SQLapiVersion": "2014-04-01"
	},
	"resources": [
		{
			"type": "Microsoft.Sql/servers/databases",
			"name": "[concat(parameters('serverName'), '/', parameters('databaseName'))]",
			"apiVersion": "[variables('SQLapiVersion')]",
			"location": "[resourceGroup().location]",
			"properties": {
				"collation": "SQL_Latin1_General_CP1_CI_AS",
				"edition": "[parameters('edition')]",
				"requestedServiceObjectiveName": "[parameters('requestedServiceObjectiveName')]",
				"maxSizeBytes": "[parameters('maxSizeBytes')]"
			},
			"tags": "[parameters('tags')]"
		}
	],
	"outputs": {
	}
}
`)
