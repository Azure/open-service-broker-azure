package mssqldr

// nolint: lll
var failoverGroupARMTemplateBytes = []byte(`
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
			"type": "Microsoft.Sql/servers/failoverGroups",
			"name": "{{ .priServerName}}/{{ .failoverGroupName }}",
			"apiVersion": "2015-05-01-preview",
			"properties": {
	      "readWriteEndpoint": {
	        "failoverPolicy": "Automatic",
	        "failoverWithDataLossGracePeriodMinutes": 60
	      },
	      "readOnlyEndpoint": {
	        "failoverPolicy": "Disabled"
	      },
	      "partnerServers": [
	        {
	          "id": "[resourceId('Microsoft.Sql/servers', '{{ .secServerName }}')]"
	        }
	      ],
	      "databases": [
	        "[resourceId('Microsoft.Sql/servers/databases', '{{ .priServerName }}', '{{ .databaseName }}')]"
	      ]
	    },
			"tags": "[parameters('tags')]"
		}
	],
	"outputs": {
	}
}
`)
