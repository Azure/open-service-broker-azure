package acr

// nolint: lll
var armTemplateBytes = []byte(`
{
    "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
    "contentVersion": "1.0.0.0",
    "parameters": {
        "registryName": {
            "type": "string",
            "metadata": {
                "description": "The name of the container registry."
            }
        },
        "location": {
            "type": "string",
            "metadata": {
                "description": "The location of the container registry. This cannot be changed after the resource is created."
            }
        },
        "registrySku": {
            "type": "string",
            "defaultValue": "Standard",
            "metadata": {
                "description": "The SKU of the container registry."
			},
			"allowedValues": [
				"Basic",
				"Standard",
				"Premium"
			  ]
        },
        "registryApiVersion": {
            "type": "string",
            "defaultValue": "2017-10-01",
            "metadata": {
                "description": "The API version of the container registry."
            }
        },
        "adminUserEnabled": {
            "type": "bool",
            "defaultValue": false,
            "metadata": {
                "description": "The value that indicates whether the admin user is enabled."
            }
        },
        "tags": {
			"type": "object"
		}
    },
    "resources": [
        {
            "name": "[parameters('registryName')]",
            "type": "Microsoft.ContainerRegistry/registries",
            "location": "[parameters('location')]",
            "apiVersion": "[parameters('registryApiVersion')]",
            "sku": {
                "name": "[parameters('registrySku')]"
            },
            "properties": {
                "adminUserEnabled": "[parameters('adminUserEnabled')]"
            }
        }
    ],
    "outputs": {
		"registryName": {
            "type": "string",
            "value": "[parameters('registryName')]"
          }
	}
}
`)
