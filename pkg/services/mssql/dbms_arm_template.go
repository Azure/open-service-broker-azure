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
		"serverName": {
			"type": "string"
		},
		"version": {
			"type": "string",
			"allowedValues": [ "2.0", "12.0" ],
			"defaultValue": "12.0"
		},
		"administratorLogin": {
			"type": "string"
		},
		"administratorLoginPassword": {
			"type": "securestring"
		},
		"tags": {
			"type": "object"
		}
	},
	"variables": {
		"SQLapiVersion": "2014-04-01"
	},
	"resources": [
		{
			"type": "Microsoft.Sql/servers",
			"name": "[parameters('serverName')]",
			"apiVersion": "[variables('SQLapiVersion')]",
			"location": "[parameters('location')]",
			"properties": {
				"administratorLogin": "[parameters('administratorLogin')]",
				"administratorLoginPassword": "[parameters('administratorLoginPassword')]",
				"version": "[parameters('version')]"
			},
			"tags": "[parameters('tags')]",
			"resources": [
				{{$count:= sub (len .firewallRules)  1}}
				{{range $i, $rule := .firewallRules}}
				{
					"type": "firewallrules",
					"name": "{{$rule.Name}}",
					"apiVersion": "[variables('SQLapiVersion')]",
					"location": "[parameters('location')]",
					"properties": {
						"startIpAddress": "{{$rule.StartIP}}",
						"endIpAddress": "{{$rule.EndIP}}"
					},
					"dependsOn": [
						"[concat('Microsoft.Sql/servers/', parameters('serverName'))]"
					]
				}{{if lt $i $count}},{{end}}
				{{end}}
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
