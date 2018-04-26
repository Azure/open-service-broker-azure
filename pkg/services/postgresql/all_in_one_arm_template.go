package postgresql

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
		"DBforPostgreSQLapiVersion": "2017-12-01"
	},
	"resources": [
		{
			"apiVersion": "[variables('DBforPostgreSQLapiVersion')]",
			"kind": "",
			"location": "[parameters('location')]",
			"name": "[parameters('serverName')]",
			"properties": {
				"version": "{{.version}}",
				"administratorLogin": "postgres",
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
			"type": "Microsoft.DBforPostgreSQL/servers",
			"tags": "[parameters('tags')]",
			"resources": [
				{{range .firewallRules}}
				{
					"type": "firewallrules",
					"apiVersion": "[variables('DBforPostgreSQLapiVersion')]",
					"dependsOn": [
						"[concat('Microsoft.DBforPostgreSQL/servers/', parameters('serverName'))]"
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
					"apiVersion": "[variables('DBforPostgreSQLapiVersion')]",
					"name": "[parameters('databaseName')]",
					"type": "databases",
					"location": "[parameters('location')]",
					"dependsOn": [
						{{range .firewallRules}}
						"[concat('Microsoft.DBforPostgreSQL/servers/', parameters('serverName'), '/firewallrules/', '{{.Name}}')]",
						{{end}}
						"[concat('Microsoft.DBforPostgreSQL/servers/', parameters('serverName'))]"
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
