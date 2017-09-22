package mysql

// nolint: lll
var armTemplateBytes = []byte(`
{
	"$schema": "http://schema.management.azure.com/schemas/2014-04-01-preview/deploymentTemplate.json#",
	"contentVersion": "1.0.0.0",
	"parameters": {
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
			"allowedValues": [ "MYSQLB50", "MYSQLB100" ]
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
			"allowedValues": [ "5.7", "5.6" ],
			"defaultValue": "5.7"
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
			"defaultValue": "255.255.255.255"
		},
		"sslEnforcement": {
			"type": "string",
			"allowedValues": [ "Enabled", "Disabled" ],
			"defaultValue": "Enabled"
		}
	},
	"variables": {
		"DBforMySQLapiVersion": "2017-04-30-preview"
	},
	"resources": [
		{
			"apiVersion": "[variables('DBforMySQLapiVersion')]",
			"kind": "",
			"location": "[resourceGroup().location]",
			"name": "[parameters('serverName')]",
			"properties": {
				"version": "[parameters('version')]",
				"administratorLogin": "azureuser",
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
			"type": "Microsoft.DBforMySQL/servers",
			"tags": {
				"heritage": "azure-service-broker"
			},
			"resources": [
				{
					"type": "firewallrules",
					"apiVersion": "[variables('DBforMySQLapiVersion')]",
					"dependsOn": [
						"[concat('Microsoft.DBforMySQL/servers/', parameters('serverName'))]"
					],
					"location": "[resourceGroup().location]",
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
					"location": "[resourceGroup().location]",
					"dependsOn": [
							"[concat('Microsoft.DBforMySQL/servers/', parameters('serverName'))]",
							"[concat('Microsoft.DBforMySQL/servers/', parameters('serverName'), '/firewallrules/', parameters('firewallRuleName'))]"
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
