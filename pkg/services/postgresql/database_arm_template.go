package postgresql

// nolint: lll
var databaseARMTemplateBytes = []byte(`
{
	"$schema": "http://schema.management.azure.com/schemas/2014-04-01-preview/deploymentTemplate.json#",
	"contentVersion": "1.0.0.0",
	"parameters": {
		"location": {
			"type": "string"
		},
		"tags": {
			"type": "object"
		}
	},
	"variables": {
		"DBforPostgreSQLapiVersion": "2017-12-01"
	},
	"resources": [
		{
			"apiVersion": "[variables('DBforPostgreSQLapiVersion')]",
			"type": "Microsoft.DBforPostgreSQL/servers/databases",
			"name": "{{.serverName}}/{{.databaseName}}",
			"properties": {}
		}
	],
	"outputs": { }
}
`)
