package keyvault

// nolint: lll
var armTemplateBytes = []byte(`
{
	"$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
	"contentVersion": "1.0.0.0",
	"parameters": {
		"tags": {
			"type": "object"
		}
	},
	"resources": [
		{
			"type": "Microsoft.KeyVault/vaults",
			"name": "{{.keyVaultName}}",
			"apiVersion": "2015-06-01",
			"location": "{{.location}}",
			"tags": "[parameters('tags')]",
			"properties": {
				"enabledForDeployment": false,
				"enabledForTemplateDeployment": false,
				"enabledForVolumeEncryption": false,
				"tenantId": "{{.tenantId}}",
				"accessPolicies": [
					{
						"objectId": "{{.objectId}}",
						"tenantId": "{{.tenantId}}",
						"permissions": {
							"keys": ["Get","List","Update","Create","Import","Delete","Recover","Backup","Restore"],
							"secrets": ["Get","List","Set","Delete","Recover","Backup","Restore"],
							"certificates": ["Get","List","Update","Create","Import","Delete","ManageContacts","ManageIssuers","GetIssuers","ListIssuers","SetIssuers","DeleteIssuers"]
						}
					}
				],
				"sku": {
					"name": "{{.vaultSku}}",
					"family": "A"
				}
			}
		}
	],
	"outputs": {
		"vaultUri": {
			"type": "string",
			"value": "[reference('{{.keyVaultName}}').vaultUri]"
		}
	}
}
`)
