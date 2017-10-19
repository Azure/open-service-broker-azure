package aci

// nolint: lll
var armTemplateBytes = []byte(`
	{
		"$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
		"contentVersion": "1.0.0.0",
		"parameters": {    
			"name": {
					"type": "string",
					"metadata": {
							"description": "Name for the container group"
					},
					"defaultValue": "acilinuxpublicipcontainergroup"
			},    
			"image": {        
					"type": "string",
					"metadata": {
							"description": "Container image to deploy. Should be of the form accountName/imagename:tag for images stored in Docker Hub or a fully qualified URI for a private registry like the Azure Container Registry."
					},
					"defaultValue": "microsoft/aci-helloworld"
			},
			"port": {
					"type": "string",
					"metadata": {
							"description": "Port to open on the container and the public IP address."
					},
					"defaultValue": "80"
			},    
			"cpuCores": {
					"type": "string",
					"metadata": {
							"description": "The number of CPU cores to allocate to the container. Must be an integer."
					},
					"defaultValue": "1.0"
			},
			"memoryInGb": {
					"type": "string",
					"metadata": {
							"description": "The amount of memory to allocate to the container in gigabytes."
					},
					"defaultValue": "1.5"
			},
			{{if .storageaccountname}}
			"sharename": {
					"type": "string"
			},
			"storageaccountname": {
					"type": "string"
			},
			"storageaccountkey": {
					"type": "securestring"
			},
			"volumename": {
					"type": "string"
			},
			"mountpoint": {
					"type": "string"
			},
			{{end}}
			{{if .imageRegistryLoginServer}}
			"imageRegistryLoginServer": {        
				"type": "string"
			},
			"imageUsername": {        
				"type": "string"
			},
			"imagePassword": {        
				"type": "string"
			},
			{{end}}
			"tags": {
				"type": "object"
			}
		},
		"variables": {},
		"resources": [
			{
							"name": "[parameters('name')]",
							"type": "Microsoft.ContainerInstance/containerGroups",
							"apiVersion": "2017-08-01-preview",
							"location": "[resourceGroup().location]",
							"properties": {
									"containers": [
											{
													"name": "[parameters('name')]",
													"properties": {
															"image": "[parameters('image')]",
															"ports": [
																	{
																			"port": "[parameters('port')]" 
																	}
															],
															{{if .storageaccountname}}
															"volumeMounts": [{
																"name": "[parameters('volumename')]",
																"mountPath": "[parameters('mountpoint')]"
															  }],
															{{end}}
															"resources": {
																	"requests": {
																			"cpu": "[parameters('cpuCores')]",
																			"memoryInGb": "[parameters('memoryInGb')]"
																	}
															}
													}
											}
									],
									"osType": "Linux",
									{{if .imageRegistryLoginServer}}
									"imageRegistryCredentials": [
									{
										"server": "[parameters('imageRegistryLoginServer')]",
										"username": "[parameters('imageUsername')]",
										"password": "[parameters('imagePassword')]"
									}
									],
									{{end}}
									{{if .storageaccountname}}
									"volumes": [{
										"name": "[parameters('volumename')]",
										"azureFile": {
											"shareName": "[parameters('sharename')]",
											"storageAccountName": "[parameters('storageaccountname')]",
											"storageAccountKey": "[parameters('storageaccountkey')]"
										}
									  }],
									{{end}}
									"ipAddress": {
											"type": "Public",
											"ports": [
													{
															"protocol": "tcp",
															"port": "[parameters('port')]"
													}
											]
									 }
							}
					}
		],
		"outputs": {
			 "containerIPv4Address":{
						"type": "string",
						"value": "[reference(resourceId('Microsoft.ContainerInstance/containerGroups/', parameters('name'))).ipAddress.ip]"
				}
		}
	}
`)
