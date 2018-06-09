// +build experimental

package mssql

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
	"resources": [
		{
			"type": "Microsoft.Sql/servers",
			"name": "{{ .serverName }}",
			"apiVersion": "2015-05-01-preview",
			"location": "[parameters('location')]",
			"properties": {
				"administratorLogin": "{{ .administratorLogin }}",
				"administratorLoginPassword": "{{ .administratorLoginPassword }}",
				"version": "{{.version}}"
			},
			"tags": "[parameters('tags')]",
			"resources": [
				{{$count:= sub (len .firewallRules)  1}}
				{{range $i, $rule := .firewallRules}}
				{
					"type": "firewallrules",
					"name": "{{$rule.Name}}",
					"apiVersion": "2014-04-01-preview",
					"location": "[parameters('location')]",
					"properties": {
						"startIpAddress": "{{$rule.StartIP}}",
						"endIpAddress": "{{$rule.EndIP}}"
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
