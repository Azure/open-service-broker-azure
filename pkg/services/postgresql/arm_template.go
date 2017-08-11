package postgresql

// nolint: lll
var armTemplateBytes = []byte(`
{
	"$schema": "http://schema.management.azure.com/schemas/2014-04-01-preview/deploymentTemplate.json#",
	"contentVersion": "1.0.0.0",
	"parameters": {
		"administratorLogin": {
			"type": "string",
			"minLength": 1,
			"maxLength": 16,
			"defaultValue": "postgres"
		},
		"administratorLoginPassword": {
			"type": "securestring"
		},
		"serverName": {
			"type": "string",
			"minLength": 2,
			"maxLength": 63
		},
		"skuCapacityDTU": {
			"type": "int",
			"allowedValues": [ 50, 100 ],
			"defaultValue": 50
		},
		"skuFamily": {
			"type": "string",
			"allowedValues": [ "SkuFamily" ],
			"defaultValue": "SkuFamily"
		},
		"skuName": {
		"type": "string",
			"allowedValues": [ "PGSQLB50", "PGSQLB100" ],
			"defaultValue": "PGSQLB50"
		},
		"skuSizeMB": {
			"type": "int",
			"minValue": 51200,
			"maxValue": 102400,
			"defaultValue": 51200
		},
		"skuTier": {
			"type": "string",
			"allowedValues": [ "Basic" ],
			"defaultValue": "Basic"
		},
		"version": {
			"type": "string",
			"allowedValues": [ "9.5", "9.6" ],
			"defaultValue": "9.6"
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
		}
	},
	"variables": {
		"DBforPostgreSQLapiVersion": "2016-02-01-privatepreview"
	},
	"resources": [
		{
			"apiVersion": "[variables('DBforPostgreSQLapiVersion')]",
			"kind": "",
			"location": "[resourceGroup().location]",
			"name": "[parameters('serverName')]",
			"properties": {
				"version": "[parameters('version')]",
				"administratorLogin": "[parameters('administratorLogin')]",
				"administratorLoginPassword": "[parameters('administratorLoginPassword')]",
				"storageMB": "[parameters('skuSizeMB')]"
			},
			"sku": {
				"name": "[parameters('skuName')]",
				"tier": "[parameters('skuTier')]",
				"capacity": "[parameters('skuCapacityDTU')]",
				"size": "[parameters('skuSizeMB')]",
				"family": "[parameters('skuFamily')]"
			},
			"type": "Microsoft.DBforPostgreSQL/servers",
			"resources": [
				{
					"type": "firewallrules",
					"apiVersion": "[variables('DBforPostgreSQLapiVersion')]",
					"dependsOn": [
						"[concat('Microsoft.DBforPostgreSQL/servers/', parameters('serverName'))]"
					],
					"location": "[resourceGroup().location]",
					"name": "[parameters('firewallRuleName')]",
					"properties": {
						"startIpAddress": "[parameters('firewallStartIpAddress')]",
						"endIpAddress": "[parameters('firewallEndIpAddress')]"
					}
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
