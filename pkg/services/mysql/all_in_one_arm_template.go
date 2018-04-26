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
		"DBforMySQLapiVersion": "2017-12-01"
	},
	"resources": [
		{
			"apiVersion": "[variables('DBforMySQLapiVersion')]",
			"kind": "",
			"location": "[parameters('location')]",
			"name": "[parameters('serverName')]",
			"properties": {
				"version": "{{.version}}",
				"administratorLogin": "azureuser",
				"administratorLoginPassword": "[parameters('administratorLoginPassword')]",
				"storageProfile" : {
					"storageMB": {{.storage}},
					{{ if .geoRedundantBackup }}
					"geoRedundantBackup" : "Enabled",
					{{ end }}
					"backupRetentionDays": {{.backupRetention}}
				},
				"sslEnforcement": "[parameters('sslEnforcement')]"
			},
			"sku": {
				"name": "{{.sku}}",
				"tier": "{{.tier}}",
				"capacity": "{{.cores}}",
				"size": "{{.storage}}",
				"family": "{{.hardwareFamily}}"
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
					"apiVersion": "[variables('DBforMySQLapiVersion')]",
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
