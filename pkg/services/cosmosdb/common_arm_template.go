package cosmosdb

// nolint: lll
var armTemplateBytes = []byte(`
{
	"$schema": "http://schema.management.azure.com/schemas/2014-04-01-preview/deploymentTemplate.json#",
	"contentVersion": "1.0.0.0",
	"parameters": {
		"tags": {
			"type": "object"
		}
	},
	"resources": [
		{
			"apiVersion": "2015-04-08",
			"kind": "{{ .kind }}",
			"type": "Microsoft.DocumentDb/databaseAccounts",
			"name": "{{ .name }}",
			"location": "{{ .location }}",
			"properties": {
				{{ if .consistencyPolicy }}
				"consistencyPolicy" : {
					"defaultConsistencyLevel" : "{{ .consistencyPolicy.defaultConsistencyLevel }}"
					{{ if .consistencyPolicy.boundedStaleness }}
						,"maxStalenessPrefix": {{ .consistencyPolicy.boundedStaleness.maxStalenessPrefix }}
					{{ end }}
					{{ if .consistencyPolicy.maxIntervalInSeconds }}
						,"maxIntervalInSeconds" : {{ .consistencyPolicy.boundedStaleness.maxIntervalInSeconds }}
					{{ end }}
				},
				{{ end }}
				"databaseAccountOfferType": "Standard",
				{{ if .ipFilters }} 
				"ipRangeFilter" : "{{ .ipFilters }}",
				{{ end }}
				{{ if .capability }}
				"capabilities": [
					{
						"name": "{{ .capability }}"
					}
				],
				{{ end }}
				"locations": [
					{{range $index,$locationInformation := .readLocations}}
					{{if $index}}
					,
					{{end}}
					{
						"failoverPriority": {{$locationInformation.Priority}},
						"locationName": "{{$locationInformation.Location}}",
						"id": "{{$locationInformation.ID}}"
					}
					{{end}}
				],
				"enableAutomaticFailover": {{ .enableAutomaticFailover}}
			},
			"tags": "[parameters('tags')]"
		}
	],
	"outputs": {
		"fullyQualifiedDomainName": {
			"type": "string",
			"value": "[reference('{{ .name }}').documentEndpoint]"
		},
		"primaryKey":{
			"type": "string",
			"value": "[listKeys(resourceId('Microsoft.DocumentDb/databaseAccounts', '{{ .name }}'), '2015-04-08').primaryMasterKey]"
		}
	}
}
`)
