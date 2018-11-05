package mssqldr

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func getDBMSPairRegisteredUpdateParamSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		SecureProperties: []string{
			"primaryAdministratorLoginPassword",
			"secondaryAdministratorLoginPassword",
		},
		PropertySchemas: map[string]service.PropertySchema{
			"primaryAdministratorLogin": &service.StringPropertySchema{
				Title: "Primary Administrator Login",
				Description: "Specifies the administrator login name" +
					" of the primary existing server",
			},
			"primaryAdministratorLoginPassword": &service.StringPropertySchema{
				Title: "Primary Administrator Login Password",
				Description: "Specifies the administrator login name" +
					" of the existing primary server",
			},
			"secondaryAdministratorLogin": &service.StringPropertySchema{
				Title: "Secondary Administrator Login",
				Description: "Specifies the administrator login name" +
					" of the secondary existing server",
			},
			"secondaryAdministratorLoginPassword": &service.StringPropertySchema{
				Title: "Secondary Administrator Login Password",
				Description: "Specifies the administrator login name" +
					" of the existing secondary server",
			},
		},
	}
}

func getDBMSPairRegisteredProvisionParamSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		RequiredProperties: []string{
			"primaryResourceGroup",
			"primaryLocation",
			"primaryServer",
			"primaryAdministratorLogin",
			"primaryAdministratorLoginPassword",
			"secondaryResourceGroup",
			"secondaryLocation",
			"secondaryServer",
			"secondaryAdministratorLogin",
			"secondaryAdministratorLoginPassword",
		},
		SecureProperties: []string{
			"primaryAdministratorLoginPassword",
			"secondaryAdministratorLoginPassword",
		},
		PropertySchemas: map[string]service.PropertySchema{
			"primaryResourceGroup": &service.StringPropertySchema{
				Title: "Primary Resource Group",
				Description: "Specifies the resource group of " +
					"the primary existing server",
			},
			"primaryLocation": &service.StringPropertySchema{
				Title: "Primary Location",
				Description: "Specifies the location of " +
					"the primary existing server",
				CustomPropertyValidator: azure.LocationValidator,
			},
			"primaryServer": &service.StringPropertySchema{
				Title:       "Primary Server",
				Description: "Specifies the name of the primary existing server",
			},
			"primaryAdministratorLogin": &service.StringPropertySchema{
				Title: "Primary Administrator Login",
				Description: "Specifies the administrator login name" +
					" of the primary existing server",
			},
			"primaryAdministratorLoginPassword": &service.StringPropertySchema{
				Title: "Primary Administrator Login Password",
				Description: "Specifies the administrator login password" +
					" of the primary existing server",
			},
			"secondaryResourceGroup": &service.StringPropertySchema{
				Title: "Secondary Resource Group",
				Description: "Specifies the resource group of " +
					"the secondary existing server",
			},
			"secondaryLocation": &service.StringPropertySchema{
				Title: "Secondary Location",
				Description: "Specifies the location of " +
					"the secondary existing server",
				CustomPropertyValidator: azure.LocationValidator,
			},
			"secondaryServer": &service.StringPropertySchema{
				Title:       "Secondary Server",
				Description: "Specifies the name of the secondary existing server",
			},
			"secondaryAdministratorLogin": &service.StringPropertySchema{
				Title: "Secondary Administrator Login",
				Description: "Specifies the administrator login name" +
					" of the secondary existing server",
			},
			"secondaryAdministratorLoginPassword": &service.StringPropertySchema{
				Title: "Secondary Administrator Login Password",
				Description: "Specifies the administrator login password" +
					" of the secondary existing server",
			},
			"tags": &service.ObjectPropertySchema{
				Title: "Tags",
				Description: "Tags to be applied to resources," +
					" specified as key/value pairs.",
				Additional: &service.StringPropertySchema{},
			},
		},
	}
}
