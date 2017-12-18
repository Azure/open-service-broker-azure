package postgresqldb

// nolint: lll
var armTemplateBytes = []byte(`
{
	"$schema": "http://schema.management.azure.com/schemas/2014-04-01-preview/deploymentTemplate.json#",
	"contentVersion": "1.0.0.0",
	"parameters": {
		"location": {
			"type": "string"
		},
		"administratorLoginPassword": {
			"type": "securestring"
		},
		"serverName": {
			"type": "string",
			"minLength": 2,
			"maxLength": 63
		},
		"skuName": {
			"type": "string",
			"allowedValues": [ "PGSQLB50", "PGSQLB100" ]
		},
		"skuTier": {
			"type": "string",
			"allowedValues": [ "Basic" ]
		},
		"skuCapacityDTU": {
			"type": "int",
			"allowedValues": [ 50, 100 ]
		},
		"skuSizeMB": {
			"type": "int",
			"minValue": 51200,
			"maxValue": 102400,
			"defaultValue": 51200
		},
		"version": {
			"type": "string",
			"allowedValues": [ "9.5", "9.6" ],
			"defaultValue": "9.6"
		},
		"databaseName": {
			"type": "string",
			"minLength": 2,
			"maxLength": 63
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
			"defaultValue": "0.0.0.0"
		},
		"sslEnforcement": {
			"type": "string",
			"allowedValues": [ "Enabled", "Disabled" ],
			"defaultValue": "Enabled"
		},
		"tags": {
			"type": "object"
		}
	},
	"variables": {
		"DBforPostgreSQLapiVersion": "2016-02-01-privatepreview"
	},
	"resources": [
		{
			"apiVersion": "[variables('DBforPostgreSQLapiVersion')]",
			"kind": "",
			"location": "[parameters('location')]",
			"name": "[parameters('serverName')]",
			"properties": {
				"version": "[parameters('version')]",
				"administratorLogin": "postgres",
				"administratorLoginPassword": "[parameters('administratorLoginPassword')]",
				"storageMB": "[parameters('skuSizeMB')]",
				"sslEnforcement": "[parameters('sslEnforcement')]"
			},
			"sku": {
				"name": "[parameters('skuName')]",
				"tier": "[parameters('skuTier')]",
				"capacity": "[parameters('skuCapacityDTU')]",
				"size": "[parameters('skuSizeMB')]"
			},
			"type": "Microsoft.DBforPostgreSQL/servers",
			"tags": "[parameters('tags')]",
			"resources": [
				{
					"type": "firewallrules",
					"apiVersion": "[variables('DBforPostgreSQLapiVersion')]",
					"dependsOn": [
						"[concat('Microsoft.DBforPostgreSQL/servers/', parameters('serverName'))]"
					],
					"location": "[parameters('location')]",
					"name": "[parameters('firewallRuleName')]",
					"properties": {
						"startIpAddress": "[parameters('firewallStartIpAddress')]",
						"endIpAddress": "[parameters('firewallEndIpAddress')]"
					}
				},
				{
					"apiVersion": "2017-04-30-preview",
					"name": "[parameters('databaseName')]",
					"type": "databases",
					"location": "[parameters('location')]",
					"dependsOn": [
						"[concat('Microsoft.DBforPostgreSQL/servers/', parameters('serverName'))]",
						"[concat('Microsoft.DBforPostgreSQL/servers/', parameters('serverName'), '/firewallrules/', parameters('firewallRuleName'))]"
					],
					"properties": {}
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
