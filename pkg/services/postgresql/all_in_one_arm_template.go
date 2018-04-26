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
		"tags": {
			"type": "object"
		}
	},
	"variables": {
		"DBforPostgreSQLapiVersion": "2017-12-01",
		"ServerName" : "{{ .serverName }}"
	},
	"resources": [
		{
			"apiVersion": "[variables('DBforPostgreSQLapiVersion')]",
			"kind": "",
			"location": "[parameters('location')]",
			"name": "[variables('ServerName')]",
			"properties": {
				"version": "{{.version}}",
				"administratorLogin": "postgres",
				"administratorLoginPassword": "{{ .administratorLoginPassword }}",
				"storageProfile" : {
					"storageMB": {{.storage}},
					{{ if .geoRedundantBackup }}
					"geoRedundantBackup" : "Enabled",
					{{ end }}
					"backupRetentionDays": {{.backupRetention}}
				},
				"sslEnforcement": "{{ .sslEnforcement }}"
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
						"[concat('Microsoft.DBforPostgreSQL/servers/', variables('ServerName'))]"
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
					"name": "{{ .databaseName }}",
					"type": "databases",
					"location": "[parameters('location')]",
					"dependsOn": [
						{{range .firewallRules}}
						"[concat('Microsoft.DBforPostgreSQL/servers/', variables('ServerName'), '/firewallrules/', '{{.Name}}')]",
						{{end}}
						"[concat('Microsoft.DBforPostgreSQL/servers/', variables('ServerName'))]"
					],
					"properties": {}
				}
			]
		}
	],
	"outputs": {
		"fullyQualifiedDomainName": {
			"type": "string",
			"value": "[reference(variables('ServerName')).fullyQualifiedDomainName]"
		}
	}
}
`)
