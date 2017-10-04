package keyvault

// nolint: lll
var armTemplateBytes = []byte(`
{
	"$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
	"contentVersion": "1.0.0.0",
	"parameters": {
		"keyVaultName": {
			"type": "string",
			"metadata": {
				"description": "Name of the Key Vault"
			}
		},
		"tenantId": {
			"type": "string",
			"metadata": {
				"description": "Tenant Id for the subscription and user assigned access to the vault."
			}
		},
		"accessPolicies": {
			"type": "array",
			"defaultValue": [],
			"metadata": {
				"description": "Access policies object"
			}
		},
		"objectId": {
			"type": "string",
			"metadata": {
				"description": "Object Id for the service principal"
			}
		},
		"vaultSku": {
			"type": "string",
			"defaultValue": "Standard",
			"allowedValues": [
				"Standard",
				"Premium"
			],
			"metadata": {
				"description": "SKU for the vault"
			}
		},
		"enabledForDeployment": {
			"type": "bool",
			"defaultValue": false,
			"metadata": {
				"description": "Specifies if the vault is enabled for VM or Service Fabric deployment"
			}
		},
		"enabledForTemplateDeployment": {
			"type": "bool",
			"defaultValue": false,
			"metadata": {
				"description": "Specifies if the vault is enabled for ARM template deployment"
			}
		},
		"enableVaultForVolumeEncryption": {
			"type": "bool",
			"defaultValue": false,
			"metadata": {
				"description": "Specifies if the vault is enabled for volume encryption"
			}
		}
	},
	"resources": [
		{
			"type": "Microsoft.KeyVault/vaults",
			"name": "[parameters('keyVaultName')]",
			"apiVersion": "2015-06-01",
			"location": "[resourceGroup().location]",
			"tags": {
				"displayName": "KeyVault"
			},
			"properties": {
				"enabledForDeployment": "[parameters('enabledForDeployment')]",
				"enabledForTemplateDeployment": "[parameters('enabledForTemplateDeployment')]",
				"enabledForVolumeEncryption": "[parameters('enableVaultForVolumeEncryption')]",
				"tenantId": "[parameters('tenantId')]",
				"accessPolicies": [{"objectId": "[parameters('objectId')]","tenantId": "[parameters('tenantId')]","permissions": {"keys": ["Get","List","Update","Create","Import","Delete","Recover","Backup","Restore"],"secrets": ["Get","List","Set","Delete","Recover","Backup","Restore"],"certificates": ["Get","List","Update","Create","Import","Delete","ManageContacts","ManageIssuers","GetIssuers","ListIssuers","SetIssuers","DeleteIssuers"]}}],
				"sku": {
					"name": "[parameters('vaultSku')]",
					"family": "A"
				}
			}
		}
	],
	"outputs": {
		"vaultUri": {
			"type": "string",
			"value": "[reference(parameters('keyVaultName')).vaultUri]"
		}
	}
}
`)
