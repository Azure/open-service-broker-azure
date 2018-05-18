package mssql

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
				{{range .firewallRules}}
				{
					"type": "firewallrules",
					"name": "{{.Name}}",
					"apiVersion": "2014-04-01-preview",
					"location": "[parameters('location')]",
					"properties": {
						"startIpAddress": "{{.StartIP}}",
						"endIpAddress": "{{.EndIP}}"
					},
					"dependsOn": [
						"Microsoft.Sql/servers/{{ $.serverName }}"
					]
				},
				{{end}}
				{
					"type": "databases",
					"name": "{{ .databaseName }}",
					"apiVersion": "2017-10-01-preview",
					"location": "[parameters('location')]",
					"properties": {
						"collation": "SQL_Latin1_General_CP1_CI_AS",
						"maxSizeBytes": "{{ .maxSizeBytes }}"
					},
					"sku" : {
						"name" : "{{ .sku }}",
						"tier" : "{{ .tier }}"
					},
					"dependsOn": [
						{{range .firewallRules}}
						"Microsoft.Sql/servers/{{ $.serverName }}/firewallrules/{{.Name}}",
						{{end}}
						"Microsoft.Sql/servers/{{ $.serverName }}"
						
					],
					"tags": "[parameters('tags')]"
				}
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
