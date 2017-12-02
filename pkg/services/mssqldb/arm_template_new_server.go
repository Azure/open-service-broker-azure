package mssqldb

// nolint: lll
var armTemplateNewServerBytes = []byte(`
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
		"firewallRuleName": {
			"type": "string",
			"minLength": 1,
			"maxLength": 128,
			"defaultValue": "AllowAll"
		},
		"firewallStartIpAddress": {
			"type": "string",
			"minLength": 1,
			"maxLength": 15,
			"defaultValue": "0.0.0.0"
		},
		"firewallEndIpAddress": {
			"type": "string",
			"minLength": 1,
			"maxLength": 15,
			"defaultValue": "255.255.255.255"
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
				{
					"type": "firewallrules",
					"name": "[parameters('firewallRuleName')]",
					"apiVersion": "[variables('SQLapiVersion')]",
					"location": "[parameters('location')]",
					"properties": {
						"startIpAddress": "[parameters('firewallStartIpAddress')]",
						"endIpAddress": "[parameters('firewallEndIpAddress')]"
					},
					"dependsOn": [
						"[concat('Microsoft.Sql/servers/', parameters('serverName'))]"
					]
				},
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
						"[concat('Microsoft.Sql/servers/', parameters('serverName'))]",
						"[concat('Microsoft.Sql/servers/', parameters('serverName'), '/firewallrules/', parameters('firewallRuleName'))]"
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
