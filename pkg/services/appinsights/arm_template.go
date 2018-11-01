package appinsights

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
			"apiVersion": "2015-05-01",
			"name": "{{.appInsightsName}}",
			"type": "Microsoft.Insights/components",
			"location": "{{.location}}",
			"properties": {
				"ApplicationId": "{{.appInsightsName}}",
				"Application_Type": "other",
				"Flow_Type": "Bluefield",
    			"Request_Source": "rest"
			},
			"tags": "[parameters('tags')]"
		}
	],
	"outputs": {
		"instrumentationKey": {
			"type": "string",
			"value": "[reference('{{.appInsightsName}}').InstrumentationKey]"
		}
	}
}
`)
