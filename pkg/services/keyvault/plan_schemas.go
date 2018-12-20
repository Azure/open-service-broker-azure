package keyvault

import (
	"github.com/Azure/open-service-broker-azure/pkg/schemas"
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
			"resourceGroup": schemas.GetResourceGroupSchema(),
			"location":      schemas.GetLocationSchema(),
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
