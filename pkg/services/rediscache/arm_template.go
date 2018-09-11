package rediscache

// nolint: lll
var armTemplateBytes = []byte(`
{
	"$schema": "http://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json",
	"contentVersion": "1.0.0.0",
	"parameters": {
		"tags": {
			"type": "object"
		}
	},
	"resources": [
		{
			"apiVersion": "2018-03-01",
			"name": "{{.serverName}}",
			"type": "Microsoft.Cache/Redis",
			"location": "{{.location}}",
			"properties": {
				{{if .redisConfiguration}}
				"redisConfiguration" : {{.redisConfiguration}},
				{{end}}
				{{if .shardCount}}
				"shardCount": "{{.shardCount}}",
				{{end}}
				{{if .subnetId}}
				"subnetId": "{{.subnetId}}",
				{{end}}
				{{if .staticIP}}
				"staticIP": "{{.staticIP}}",
				{{end}}
				"enableNonSslPort": {{.enableNonSslPort}},
				"sku": {
					"capacity": "{{.redisCacheCapacity}}",
					"family": "{{.redisCacheFamily}}",
					"name": "{{.redisCacheSKU}}"
				}
			},
			"tags": "[parameters('tags')]"
		}
	],
	"outputs": {
		"primaryKey": {
				"type": "string",
				"value": "[listKeys(resourceId('Microsoft.Cache/Redis', '{{.serverName}}'), '2018-03-01').primaryKey]"
		},
		"fullyQualifiedDomainName": {
				"type": "string",
				"value": "[reference('{{.serverName}}').hostName]"
		}
	}
}
`)
