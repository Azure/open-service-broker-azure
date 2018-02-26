package sqldb

// nolint: lll
var armTemplateDBMSOnlyBytes = []byte(`
{
	"$schema": "http://schema.management.azure.com/schemas/2014-04-01-preview/deploymentTemplate.json#",
	"contentVersion": "1.0.0.0",
	"parameters": {
		"location": {
			"type": "string"
		},
		"serverName": {
			"type": "string"
		},
		"version": {
			"type": "string",
			"allowedValues": [ "2.0", "12.0" ],
			"defaultValue": "12.0"
		},
		"administratorLogin": {
			"type": "string"
		},
		"administratorLoginPassword": {
			"type": "securestring"
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
			"type": "Microsoft.Sql/servers",
			"name": "[parameters('serverName')]",
			"apiVersion": "[variables('SQLapiVersion')]",
			"location": "[parameters('location')]",
			"properties": {
				"administratorLogin": "[parameters('administratorLogin')]",
				"administratorLoginPassword": "[parameters('administratorLoginPassword')]",
				"version": "[parameters('version')]"
			},
			"tags": "[parameters('tags')]",
			"resources": [
				{{range .firewallRules}}
				{
					"type": "firewallrules",
					"name": "{{.Name}}",
					"apiVersion": "[variables('SQLapiVersion')]",
					"location": "[parameters('location')]",
					"properties": {
						"startIpAddress": "{{.StartIP}}",
						"endIpAddress": "{{.EndIP}}"
					},
					"dependsOn": [
						"[concat('Microsoft.Sql/servers/', parameters('serverName'))]"
					]
				},
				{{end}}
				{
					"type": "databases",
					"name": "[parameters('databaseName')]",
					"apiVersion": "[variables('SQLapiVersion')]",
					"location": "[parameters('location')]",
					"properties": {
						"collation": "SQL_Latin1_General_CP1_CI_AS",
						"edition": "[parameters('edition')]",
						"requestedServiceObjectiveName": "[parameters('requestedServiceObjectiveName')]",
						"maxSizeBytes": "[parameters('maxSizeBytes')]"
					},
					"dependsOn": [
						{{range .firewallRules}}
						"[concat('Microsoft.Sql/servers/', parameters('serverName'), '/firewallrules/', '{{.Name}}')]",
						{{end}}
						"[concat('Microsoft.Sql/servers/', parameters('serverName'))]"
						
					],
					"tags": "[parameters('tags')]"
				}
			]
		}
	],
	"outputs": {
		"fullyQualifiedDomainName": {
			"type": "string",
			"value": "[reference(parameters('serverName')).fullyQualifiedDomainName]"
		}
	}
}
`)
