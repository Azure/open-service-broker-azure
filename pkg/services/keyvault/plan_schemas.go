package keyvault

import (
	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	s *serviceManager,
) getProvisionParametersSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		RequiredProperties: []string{
			"location",
			"resourceGroup",
			"objectId",
			"clientId",
			"clientSecret",
		},
		PropertySchemas: map[string]service.PropertySchema{
			"location": &service.StringPropertySchema{
				Title: "Location",
				Description: "The Azure region in which to provision" +
					" applicable resources.",
				CustomPropertyValidator: azure.LocationValidator,
			},
			"resourceGroup": &service.StringPropertySchema{
				Title: "Resource group",
				Description: "The (new or existing) resource group with which" +
					" to associate new resources.",
			},
			"objectId": &service.StringPropertySchema{
				Title: "Object ID",
				Description: "Object ID for an existing service principal, " +
					"which will be granted access to the new vault.",
			},
			"clientId": &service.StringPropertySchema{
				Title: "Client ID",
				Description: "Client ID (username) for an existing service principal," +
					"which will be granted access to the new vault.",
			},
			"clientSecret": &service.StringPropertySchema{
				Title: "Client secret",
				Description: "Client secret (password) for an existing service " +
					"principal, which will be granted access to the new vault.",
			},
			"tags": &service.ObjectPropertySchema{
				Title: "Tags",
				Description: "Tags to be applied to new resources," +
					" specified as key/value pairs.",
				Additional: &service.StringPropertySchema{},
			},
		},
	}
}
