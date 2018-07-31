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
				Description: "The Azure region in which to provision" +
					" applicable resources.",
				CustomPropertyValidator: azure.LocationValidator,
			},
			"resourceGroup": &service.StringPropertySchema{
				Description: "The (new or existing) resource group with which" +
					" to associate new resources.",
			},
			"objectId": &service.StringPropertySchema{
				Description: "Object ID for an existing service principal, " +
					"which will be granted access to the new vault.",
			},
			"clientId": &service.StringPropertySchema{
				Description: "Client ID (username) for an existing service principal," +
					"which will be granted access to the new vault.",
			},
			"clientSecret": &service.StringPropertySchema{
				Description: "Client secret (password) for an existing service " +
					"principal, which will be granted access to the new vault.",
			},
		},
	}
}
