package keyvault

import "github.com/Azure/open-service-broker-azure/pkg/service"

type provisioningParameters struct {
	ObjectID string `json:"objectId"`
	ClientID string `json:"clientId"`
}

func (
	s *serviceManager,
) getProvisionParametersSchema() map[string]service.ParameterSchema {

	p := map[string]service.ParameterSchema{}

	objectIDSchema := service.NewSimpleParameterSchema(
		"string",
		"Object ID for an existing service principal, "+
			"which will be granted access to the new vault.",
	)
	objectIDSchema.SetRequired(true)
	p["objectId"] = objectIDSchema

	clientIDSchema := service.NewSimpleParameterSchema(
		"string",
		"Client ID (username) for an existing service principal,"+
			"which will be granted access to the new vault.",
	)
	clientIDSchema.SetRequired(true)
	p["clientId"] = clientIDSchema

	clientSecretSchema := service.NewSimpleParameterSchema(
		"string",
		"Client secret (password) for an existing service "+
			"principal, which will be granted access to the new vault.",
	)
	clientSecretSchema.SetRequired(true)
	p["clientSecret"] = clientSecretSchema

	return p
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
	cpp service.CombinedProvisioningParameters,
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
	params service.CombinedBindingParameters,
) (service.BindingParameters, service.SecureBindingParameters, error) {
	return nil, nil, nil
}
