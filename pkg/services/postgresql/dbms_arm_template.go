package postgresql

// nolint: lll
var dbmsARMTemplateBytes = []byte(`
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
		"DBforPostgreSQLapiVersion": "2017-12-01"
	},
	"resources": [
		{
			"apiVersion": "[variables('DBforPostgreSQLapiVersion')]",
			"kind": "",
			"location": "[parameters('location')]",
			"name": "{{ .serverName }}",
			"properties": {
				"version": "{{.version}}",
				"administratorLogin": "postgres",
				"administratorLoginPassword": "{{ .administratorLoginPassword }}",
				"storageProfile": {
					"storageMB": {{.storage}},
					{{ if .geoRedundantBackup }}
					"geoRedundantBackup": "Enabled",
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
				{{$count:= sub (len .firewallRules)  1}}
				{{range $i, $rule := .firewallRules}}
				{
					"type": "firewallrules",
					"apiVersion": "[variables('DBforPostgreSQLapiVersion')]",
					"dependsOn": [
						"Microsoft.DBforPostgreSQL/servers/{{ $.serverName }}"
					],
					"location": "[parameters('location')]",
					"name": "{{$rule.Name}}",
					"properties": {
						"startIpAddress": "{{$rule.StartIP}}",
						"endIpAddress": "{{$rule.EndIP}}"
					}
				}{{if lt $i $count}},{{end}}
				{{end}}
			]
		}
	],
	"outputs": {
		"fullyQualifiedDomainName": {
			"type": "string",
			"value": "[reference('{{ .serverName }}').fullyQualifiedDomainName]"
		}
	}
}
`)
