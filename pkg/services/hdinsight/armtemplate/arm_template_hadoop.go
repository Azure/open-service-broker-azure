package armtemplate

// Hadoop is the bytes of Azure ARM Template to deploy a HDInisght whose type
// is Hadoop. It is based on the template on Azure Portal for quick creation.
// nolint: lll
func Hadoop() []byte {
	return []byte(`
{
	"$schema": "http://schema.management.azure.com/schemas/2014-04-01-preview/deploymentTemplate.json#",
	"contentVersion": "0.9.0.0",
	"parameters": {
		"clusterName": {
			"type": "String",
			"metadata": {
				"description": "The name of the HDInsight cluster to create."
			}
		},
		"clusterLoginUserName": {
			"defaultValue": "admin",
			"type": "String",
			"metadata": {
				"description": "These credentials can be used to submit jobs to the cluster and to log into cluster dashboards."
			}
		},
		"clusterLoginPassword": {
			"type": "SecureString",
			"metadata": {
				"description": "The password must be at least 10 characters in length and must contain at least one digit, one non-alphanumeric character, and one upper or lower case letter."
			}
		},
		"clusterVersion": {
			"defaultValue": "3.6",
			"type": "String",
			"metadata": {
				"description": "HDInsight cluster version."
			}
		},
		"clusterWorkerNodeCount": {
			"defaultValue": 4,
			"type": "Int",
			"metadata": {
				"description": "The number of nodes in the HDInsight cluster."
			}
		},
		"clusterKind": {
			"defaultValue": "hadoop",
			"type": "String",
			"metadata": {
				"description": "The type of the HDInsight cluster to create."
			}
		},
		"sshUserName": {
			"defaultValue": "sshuser",
			"type": "String",
			"metadata": {
				"description": "These credentials can be used to remotely access the cluster."
			}
		},
		"sshPassword": {
			"type": "SecureString",
			"metadata": {
				"description": "The password must be at least 10 characters in length and must contain at least one digit, one non-alphanumeric character, and one upper or lower case letter."
			}
		},
		"storageAccountName": {
			"type": "String"
		},
		"blobStorageContainerName": {
			"type": "String"
		},
		"blobStorageEndpoint": {
			"type": "String"
		},
		"tags": {
			"type": "object"
		}
	},
	"variables": {
		"HDInsightApiVersion": "2015-03-01-preview",
		"StorageApiVersion": "2015-05-01-preview"
	},
	"resources": [
		{
			"type": "Microsoft.HDInsight/clusters",
			"name": "[parameters('clusterName')]",
			"apiVersion": "[variables('HDInsightApiVersion')]",
			"location": "[resourceGroup().location]",
			"properties": {
				"clusterVersion": "[parameters('clusterVersion')]",
				"osType": "Linux",
				"tier": "standard",
				"clusterDefinition": {
					"kind": "[parameters('clusterKind')]",
					"configurations": {
						"gateway": {
							"restAuthCredential.isEnabled": true,
							"restAuthCredential.username": "[parameters('clusterLoginUserName')]",
							"restAuthCredential.password": "[parameters('clusterLoginPassword')]"
						}
					}
				},
				"storageProfile": {
					"storageaccounts": [
						{
							"name": "[parameters('blobStorageEndpoint')]",
							"isDefault": true,
							"container": "[parameters('blobStorageContainerName')]",
							"key": "[listKeys(resourceId('Microsoft.Storage/storageAccounts', parameters('storageAccountName')), variables('StorageApiVersion')).key1]"
						}
					]
				},
				"computeProfile": {
					"roles": [
						{
							"name": "headnode",
							"targetInstanceCount": 2,
							"hardwareProfile": {
								"vmSize": "Standard_D12_V2"
							},
							"osProfile": {
								"linuxOperatingSystemProfile": {
									"username": "[parameters('sshUserName')]",
									"password": "[parameters('sshPassword')]"
								}
							},
							"virtualNetworkProfile": null,
							"scriptActions": []
						},
						{
							"name": "workernode",
							"targetInstanceCount": "[parameters('clusterWorkerNodeCount')]",
							"hardwareProfile": {
								"vmSize": "Standard_D4_V2"
							},
							"osProfile": {
								"linuxOperatingSystemProfile": {
									"username": "[parameters('sshUserName')]",
									"password": "[parameters('sshPassword')]"
								}
							},
							"virtualNetworkProfile": null,
							"scriptActions": []
						},
						{
							"name": "zookeepernode",
							"targetInstanceCount": 3,
							"hardwareProfile": {
								"vmSize": "Small"
							},
							"osProfile": {
								"linuxOperatingSystemProfile": {
									"username": "[parameters('sshUserName')]",
									"password": "[parameters('sshPassword')]"
								}
							},
							"virtualNetworkProfile": null,
							"scriptActions": []
						}
					]
				},
				"tags": "[parameters('tags')]"
			},
			"dependsOn": [
				"[concat('Microsoft.Storage/storageAccounts/', parameters('storageAccountName'))]"
			]
		},
		{
			"type": "Microsoft.Storage/storageAccounts",
			"name": "[parameters('storageAccountName')]",
			"apiVersion": "[variables('StorageApiVersion')]",
			"location": "[resourceGroup().location]",
			"properties": {
				"accountType": "Standard_LRS"
			}
		}
	],
	"outputs": {
		"storageAccountKey": {
			"type": "string",
			"value": "[listKeys(resourceId('Microsoft.Storage/storageAccounts', parameters('storageAccountName')), variables('StorageApiVersion')).key1]"
		}
	}
}
`)
}
