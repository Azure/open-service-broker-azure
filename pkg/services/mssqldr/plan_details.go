package mssqldr

import (
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
				Description: "Specifies the administrator login name" +
					" of the primary existing server",
			},
			"primaryAdministratorLoginPassword": &service.StringPropertySchema{
				Description: "Specifies the administrator login name" +
					" of the existing primary server",
			},
			"secondaryAdministratorLogin": &service.StringPropertySchema{
				Description: "Specifies the administrator login name" +
					" of the secondary existing server",
			},
			"secondaryAdministratorLoginPassword": &service.StringPropertySchema{
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
				Description: "Specifies the resource group of " +
					"the primary existing server",
			},
			"primaryLocation": &service.StringPropertySchema{
				Description: "Specifies the location of the primary existing server",
			},
			"primaryServer": &service.StringPropertySchema{
				Description: "Specifies the name of the primary existing server",
			},
			"primaryAdministratorLogin": &service.StringPropertySchema{
				Description: "Specifies the administrator login name" +
					" of the primary existing server",
			},
			"primaryAdministratorLoginPassword": &service.StringPropertySchema{
				Description: "Specifies the administrator login password" +
					" of the primary existing server",
			},
			"secondaryResourceGroup": &service.StringPropertySchema{
				Description: "Specifies the resource group of " +
					"the secondary existing server",
			},
			"secondaryLocation": &service.StringPropertySchema{
				Description: "Specifies the location of the secondary existing server",
			},
			"secondaryServer": &service.StringPropertySchema{
				Description: "Specifies the name of the secondary existing server",
			},
			"secondaryAdministratorLogin": &service.StringPropertySchema{
				Description: "Specifies the administrator login name" +
					" of the secondary existing server",
			},
			"secondaryAdministratorLoginPassword": &service.StringPropertySchema{
				Description: "Specifies the administrator login password" +
					" of the secondary existing server",
			},
			"tags": &service.ObjectPropertySchema{
				Description: "Tags to be applied to resources," +
					" specified as key/value pairs.",
				Additional: &service.StringPropertySchema{},
			},
		},
	}
}