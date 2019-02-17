package postgresql

// nolint: lll
var dbmsARMTemplateBytes = []byte(`
{
	"$schema": "http://schema.management.azure.com/schemas/2014-04-01-preview/deploymentTemplate.json#",
	"contentVersion": "1.0.0.0",
	"parameters": {
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
			"location": "{{.location}}",
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
				"family": "Gen5"
			},
			"type": "Microsoft.DBforPostgreSQL/servers",
			"tags": "[parameters('tags')]",
			"resources": [
				{{ $root := . }}
				{{$firewallRulesCount := sub (len .firewallRules)  1}}
				{{$virtualNetworkRulesCount := sub (len .virtualNetworkRules) 1 }}
				{{range $i, $rule := .firewallRules}}
				{
					"type": "firewallrules",
					"apiVersion": "[variables('DBforPostgreSQLapiVersion')]",
					"dependsOn": [
						"Microsoft.DBforPostgreSQL/servers/{{ $.serverName }}"
					],
					"location": "{{$root.location}}",
					"name": "{{$rule.name}}",
					"properties": {
						"startIpAddress": "{{$rule.startIPAddress}}",
						"endIpAddress": "{{$rule.endIPAddress}}"
					}
				}{{ if or (lt $i $firewallRulesCount) (gt $virtualNetworkRulesCount -1) }},{{end}}
				{{end}}	
				{{range $i, $rule := .virtualNetworkRules}}
				{
					"type": "virtualNetworkRules",
					"apiVersion": "[variables('DBforPostgreSQLapiVersion')]",
					"dependsOn": [
						"Microsoft.DBforPostgreSQL/servers/{{ $.serverName }}"
					],
					"location": "{{$root.location}}",
					"name": "{{$rule.name}}",
					"properties": {
						"ignoreMissingVnetServiceEndpoint": true,
        				"virtualNetworkSubnetId": "{{$rule.subnetId}}"
					}
				}{{if lt $i $virtualNetworkRulesCount}},{{end}}
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
