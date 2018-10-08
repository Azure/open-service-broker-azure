package mssql

// nolint: lll
var databaseFeARMTemplateBytes = []byte(`
{
	"$schema": "http://schema.management.azure.com/schemas/2014-04-01-preview/deploymentTemplate.json#",
	"contentVersion": "1.0.0.0",
	"parameters": {
		"tags": {
			"type": "object"
		}
	},
	"resources": [
		{
			"type": "Microsoft.Sql/servers/databases",
			"name": "{{ .serverName}}/{{ .databaseName }}",
			"apiVersion": "2017-10-01-preview",
			"location": "{{.location}}",
			"tags": "[parameters('tags')]"
		}
	],
	"outputs": {
	}
}
`)
