package mssql

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
	"resources": [
		{
			"type": "Microsoft.Sql/servers",
			"name": "{{ .serverName }}",
			"apiVersion": "2015-05-01-preview",
			"location": "{{.location}}",
			"properties": {
				"administratorLogin": "{{ .administratorLogin }}",
				"administratorLoginPassword": "{{ .administratorLoginPassword }}",
				"version": "{{.version}}"
			},
			"tags": "[parameters('tags')]",
			"resources": [
				{{ $root := . }}
				{{$count:= sub (len .firewallRules)  1}}
				{{range $i, $rule := .firewallRules}}
				{
					"type": "firewallrules",
					"name": "{{$rule.name}}",
					"apiVersion": "2014-04-01-preview",
					"location": "{{$root.location}}",
					"properties": {
						"startIpAddress": "{{$rule.startIPAddress}}",
						"endIpAddress": "{{$rule.endIPAddress}}"
					},
					"dependsOn": [
						"Microsoft.Sql/servers/{{$.serverName}}"
					]
				}{{if lt $i $count}},{{end}}
				{{end}}
			]
		}
	],
	"outputs": {
		"fullyQualifiedDomainName": {
			"type": "string",
			"value": "[reference('{{ .serverName}}').fullyQualifiedDomainName]"
		}
	}
}
`)
