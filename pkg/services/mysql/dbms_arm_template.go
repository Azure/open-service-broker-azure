package mysql

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
			"DBforMySQLapiVersion": "2017-12-01",
			"ServerName" : "{{ .serverName }}"
		},
		"resources": [
			{
				"apiVersion": "[variables('DBforMySQLapiVersion')]",
				"kind": "",
				"location": "[parameters('location')]",
				"name": "[variables('ServerName')]",
				"properties": {
					"version": "{{.version}}",
					"administratorLogin": "azureuser",
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
				"type": "Microsoft.DBforMySQL/servers",
				"tags": "[parameters('tags')]",
				"resources": [
					{{$count:= sub (len .firewallRules)  1}}
					{{range $i, $rule := .firewallRules}}
					{
						"type": "firewallrules",
						"apiVersion": "[variables('DBforMySQLapiVersion')]",
						"dependsOn": [
							"[concat('Microsoft.DBforMySQL/servers/', variables('ServerName'))]"
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
				"value": "[reference(variables('ServerName')).fullyQualifiedDomainName]"
			}
		}
	}
	`)
