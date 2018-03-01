package mysql

// nolint: lll
var databaseARMTemplateBytes = []byte(`
	{
		"$schema": "http://schema.management.azure.com/schemas/2014-04-01-preview/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"parameters": {
			"version": {
				"type": "string",
				"allowedValues": [ "5.7", "5.6" ],
				"defaultValue": "5.7"
			},
			"location": {
				"type": "string"
			},
			"serverName": {
				"type": "string",
				"minLength": 2,
				"maxLength": 63
			},
			"databaseName": {
				"type": "string",
				"minLength": 2,
				"maxLength": 63
			},
			"tags": {
				"type": "object"
			}
		},
		"variables": {
			"DBforMySQLapiVersion": "2017-04-30-preview"
		},
		"resources": [
			{
				"apiVersion": "2017-04-30-preview",
				"type": "Microsoft.DBforMySQL/servers/databases",
				"name": "[concat(parameters('serverName'), '/', parameters('databaseName'))]",
				"properties": {}
			}
		],
		"outputs": {
		}
	}
	`)
