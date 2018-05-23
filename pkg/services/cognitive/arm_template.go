package cognitive

// nolint: lll
var armTemplateBytes = []byte(`
{
    "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
    "contentVersion": "1.0.0.0",
    "parameters": {
        "location": {
            "type": "string"
        },
        "name": {
            "defaultValue": "myTwitterSentimentAPI",
            "type": "String"
        },
        "tier": {
            "type": "String"
        },
        "tags": {
            "type": "object"
        }
    },
    "variables": {
        "cognitiveservicesid": "[concat(resourceGroup().id,'/providers/','Microsoft.CognitiveServices/accounts/', parameters('name'))]"
    },
    "resources": [{
        "type": "Microsoft.CognitiveServices/accounts",
        "sku": {
            "name": "[parameters('tier')]"
        },
        "kind": "TextAnalytics",
        "name": "[parameters('name')]",
        "apiVersion": "2016-02-01-preview",
        "location": "[parameters('location')]",
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
        }
    }
}
`)
