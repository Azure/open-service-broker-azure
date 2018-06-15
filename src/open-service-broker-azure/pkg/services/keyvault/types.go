// +build experimental

package keyvault

import "open-service-broker-azure/pkg/service"

type provisioningParameters struct {
	ObjectID string `json:"objectId"`
	ClientID string `json:"clientId"`
}

func (
	s *serviceManager,
) getProvisionParametersSchema() service.InputParametersSchema {
	return service.InputParametersSchema{
		RequiredProperties: []string{
			"objectId",
			"clientId",
			"clientSecret",
		},
		PropertySchemas: map[string]service.PropertySchema{
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

type secureProvisioningParameters struct {
	ClientSecret string `json:"clientSecret"`
}

type instanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	KeyVaultName      string `json:"keyVaultName"`
	VaultURI          string `json:"vaultUri"`
	ClientID          string `json:"clientId"`
}

type credentials struct {
	VaultURI     string `json:"vaultUri"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

func (s *serviceManager) SplitProvisioningParameters(
	cpp map[string]interface{},
) (
	service.ProvisioningParameters,
	service.SecureProvisioningParameters,
	error,
) {
	pp := provisioningParameters{}
	if err := service.GetStructFromMap(cpp, &pp); err != nil {
		return nil, nil, err
	}
	spp := secureProvisioningParameters{}
	if err := service.GetStructFromMap(cpp, &spp); err != nil {
		return nil, nil, err
	}
	ppMap, err := service.GetMapFromStruct(pp)
	if err != nil {
		return nil, nil, err
	}
	sppMap, err := service.GetMapFromStruct(spp)
	return ppMap, sppMap, err
}

func (s *serviceManager) SplitBindingParameters(
	params map[string]interface{},
) (service.BindingParameters, service.SecureBindingParameters, error) {
	return nil, nil, nil
}
