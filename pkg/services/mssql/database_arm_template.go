package mssql

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
	"resources": [
		{
			"type": "Microsoft.Sql/servers/databases",
			"name": "{{ .serverName}}/{{ .databaseName }}",
			"apiVersion": "2017-10-01-preview",
			"location": "[parameters('location')]",
			"properties": {
				"collation": "SQL_Latin1_General_CP1_CI_AS",
				"maxSizeBytes": "{{ .maxSizeBytes }}"
			},
			"sku" : {
				"name" : "{{ .sku }}",
				"tier" : "{{ .tier }}"
			},
			"tags": "[parameters('tags')]"
		}
	],
	"outputs": {
	}
}
`)
