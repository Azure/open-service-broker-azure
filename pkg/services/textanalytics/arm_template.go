package textanalytics

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
    "variables": {
        "cognitiveservicesid": "[concat(resourceGroup().id,'/providers/','Microsoft.CognitiveServices/accounts/','{{ .name }}')]"
    },
    "resources": [{
        "type": "Microsoft.CognitiveServices/accounts",
        "sku": {
            "name": "{{ .tier }}"
        },
        "kind": "TextAnalytics",
        "name": "{{ .name }}",
        "apiVersion": "2016-02-01-preview",
        "location": "{{ .location }}",
        "scale": null,
        "properties": {},
        "dependsOn": []
    }],
    "outputs": {
        "cognitivekey": {
            "type": "string",
            "value": "[listKeys(variables('cognitiveservicesid'),'2016-02-01-preview').key1]"
        },
        "endpoint": {
            "type": "string",
            "value": "[reference(variables('cognitiveservicesid'),'2016-02-01-preview').endpoint]"
        },
        "name": {
            "type": "string",
            "value": "{{ .name }}"
        }
    }
}
`)
