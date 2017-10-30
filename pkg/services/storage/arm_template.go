package storage

// nolint: lll
var armTemplateBytesGeneralPurposeStorage = []byte(`
{
	"$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
	"contentVersion": "1.0.0.0",
	"parameters": {
		"location": {
			"type": "string"
		},
		"name": {
			"type": "string"
		},
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
		"blobEncryptionEnabled": {
			"type": "bool",
			"defaultValue": false
		},
		"tags": {
			"type": "object"
		}
	},
	"resources": [
		{
			"type": "Microsoft.Storage/storageAccounts",
			"name": "[parameters('name')]",
			"apiVersion": "2017-06-01",
			"location": "[parameters('location')]",
			"sku": {
				"name": "[parameters('accountType')]"
			},
			"kind": "Storage",
			"properties": {
				"encryption": {
					"keySource": "Microsoft.Storage",
					"services": {
						"blob": {
							"enabled": "[parameters('blobEncryptionEnabled')]"
						}
					}
				}
			},
			"tags": "[parameters('tags')]"
		}
	],
	"outputs": {
		"accessKey": {
			"type": "string",
			"value": "[listKeys(resourceId('Microsoft.Storage/storageAccounts', parameters('name')), '2015-06-15').key1]"
		}
	}
}
`)

// nolint: lll
var armTemplateBytesBlobStorage = []byte(`
{
	"$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
	"contentVersion": "1.0.0.0",
	"parameters": {
		"location": {
			"type": "string"
		},
		"name": {
			"type": "string"
		},
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
		"accessTier": {
			"type": "string",
			"defaultValue": "Hot",
			"allowedValues": [
				"Cool",
				"Hot"
			]
		},
		"blobEncryptionEnabled": {
			"type": "bool",
			"defaultValue": false
		},
		"tags": {
			"type": "object"
		}
	},
	"resources": [
		{
			"type": "Microsoft.Storage/storageAccounts",
			"name": "[parameters('name')]",
			"apiVersion": "2017-06-01",
			"location": "[parameters('location')]",
			"sku": {
				"name": "[parameters('accountType')]"
			},
			"kind": "BlobStorage",
			"properties": {
				"accessTier": "[parameters('accessTier')]",
				"encryption": {
					"keySource": "Microsoft.Storage",
					"services": {
						"blob": {
							"enabled": "[parameters('blobEncryptionEnabled')]"
						}
					}
				}
			},
			"tags": "[parameters('tags')]"
		}
	],
	"outputs": {
		"accessKey": {
			"type": "string",
			"value": "[listKeys(resourceId('Microsoft.Storage/storageAccounts', parameters('name')), '2015-06-15').key1]"
		}
	}
}
`)
