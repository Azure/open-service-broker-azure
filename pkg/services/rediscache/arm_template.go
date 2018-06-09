// +build experimental

package rediscache

// nolint: lll
var armTemplateBytes = []byte(`
{
	"$schema": "http://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json",
	"contentVersion": "1.0.0.0",
	"parameters": {
		"location": {
			"type": "string"
		},
		"serverName": {
			"type": "string",
			"metadata": {
				"description": "The name of the Azure Redis Cache to create."
			}
		},
		"redisCacheSKU": {
			"type": "string",
			"allowedValues": [
				"Basic",
				"Standard",
				"Premium"
			],
			"metadata": {
				"description": "The pricing tier of the new Azure Redis Cache."
			}
		},
		"redisCacheFamily": {
			"type": "string",
			"metadata": {
				"description": "The family for the sku."
			},
			"allowedValues": [
				"C",
				"P"
			]
		},
		"redisCacheCapacity": {
			"type": "int",
			"allowedValues": [
				0,
				1,
				2,
				3,
				4,
				5,
				6
			],
			"metadata": {
				"description": "The size of the new Azure Redis Cache instance. "
			}
		},
		"diagnosticsStatus": {
			"type": "string",
			"defaultValue": "OFF",
			"allowedValues": [
				"ON",
				"OFF"
			],
			"metadata": {
				"description": "A value that indicates whether diagnostices is enabled. Use ON or OFF."
			}
		},
		"tags": {
			"type": "object"
		}
	},
	"resources": [
		{
			"apiVersion": "2015-08-01",
			"name": "[parameters('serverName')]",
			"type": "Microsoft.Cache/Redis",
			"location": "[parameters('location')]",
			"properties": {
				"enableNonSslPort": true,
				"sku": {
					"capacity": "[parameters('redisCacheCapacity')]",
					"family": "[parameters('redisCacheFamily')]",
					"name": "[parameters('redisCacheSKU')]"
				}
			},
			"tags": "[parameters('tags')]",
			"resources": [
				{
					"apiVersion": "2015-07-01",
					"type": "Microsoft.Cache/redis/providers/diagnosticsettings",
					"name": "[concat(parameters('serverName'), '/Microsoft.Insights/service')]",
					"location": "[parameters('location')]",
					"dependsOn": [
						"[concat('Microsoft.Cache/Redis/', parameters('serverName'))]"
					],
					"properties": {
						"status": "[parameters('diagnosticsStatus')]"
					}
				}
			]
		}
	],
	"outputs": {
		"primaryKey": {
				"type": "string",
				"value": "[listKeys(resourceId('Microsoft.Cache/Redis', parameters('serverName')), '2015-08-01').primaryKey]"
		},
		"fullyQualifiedDomainName": {
				"type": "string",
				"value": "[reference(parameters('serverName')).hostName]"
		}
	}
}
`)
