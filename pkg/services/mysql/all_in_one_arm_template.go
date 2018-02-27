package mysql

// nolint: lll
var allInOneARMTemplateBytes = []byte(`
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
			"allowedValues": [ "MYSQLB50", "MYSQLB100", "MYSQLS100","MYSQLS200", "MYSQLS400", "MYSQLS800" ]
		},
		"skuTier": {
			"type": "string",
			"allowedValues": [ "Basic", "Standard"]
		},
		"skuCapacityDTU": {
			"type": "int",
			"allowedValues": [ 50, 100, 200, 400, 800 ]
		},
		"skuSizeMB": {
			"type": "int",
			"minValue": 51200,
			"maxValue": 128000
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
		"DBforMySQLapiVersion": "2017-04-30-preview"
	},
	"resources": [
		{
			"apiVersion": "[variables('DBforMySQLapiVersion')]",
			"kind": "",
			"location": "[parameters('location')]",
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
			"tags": "[parameters('tags')]",
			"resources": [
				{{range .firewallRules}}
				{
					"type": "firewallrules",
					"apiVersion": "[variables('DBforMySQLapiVersion')]",
					"dependsOn": [
						"[concat('Microsoft.DBforMySQL/servers/', parameters('serverName'))]"
					],
					"location": "[parameters('location')]",
					"name": "{{.Name}}",
					"properties": {
						"startIpAddress": "{{.StartIP}}",
						"endIpAddress": "{{.EndIP}}"
					}
				},
				{{end}}
				{
					"apiVersion": "2017-04-30-preview",
					"name": "[parameters('databaseName')]",
					"type": "databases",
					"location": "[parameters('location')]",
					"dependsOn": [
						{{range .firewallRules}}
						"[concat('Microsoft.DBforMySQL/servers/', parameters('serverName'), '/firewallrules/', '{{.Name}}')]",
						{{end}}
							"[concat('Microsoft.DBforMySQL/servers/', parameters('serverName'))]"
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
