package storage

// nolint: lll
var armTemplateBytes = []byte(`
{
	"$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
	"contentVersion": "1.0.0.0",
	"parameters": {
		"accountType": {
			"type": "string",
			"defaultValue": "Standard_LRS",
			"allowedValues": [
				"Standard_LRS",
				"Standard_GRS",
				"Standard_RAGRS",
				"Standard_ZRS",
				"Premium_LRS"
			]
		},
		"tags": {
			"type": "object"
		}
	},
	"resources": [
		{
			"type": "Microsoft.Storage/storageAccounts",
			"name": "{{ .name }}",
			"apiVersion": "2017-06-01",
			"location": "{{ .location }}",
			"sku": {
				"name": "[parameters('accountType')]"
			},
			"kind": "{{.kind}}",
			"properties": {
				{{ if .accessTier }}
				"accessTier": "{{.accessTier}}",
				{{ end }}
				"supportsHttpsTrafficOnly": true
			},
			"tags": "[parameters('tags')]"
		}
	],
	"outputs": {
		"accessKey": {
			"type": "string",
			"value": "[listKeys(resourceId('Microsoft.Storage/storageAccounts', '{{ .name }}'), '2015-06-15').key1]"
		}
	}
}
`)
