package keyvault

import "github.com/Azure/open-service-broker-azure/pkg/service"

type provisioningParameters struct {
	ObjectID string `json:"objectId"`
	ClientID string `json:"clientId"`
}

// GetProvisionParametersSchema generates the schema for instance provisioning parameters
func GetProvisionParametersSchema() *service.ParametersSchema {

	p := service.GetCommonProvisionParametersSchema()

	p.Properties["objectId"] = service.Parameter{
		Type: "string",
		Description: "Object ID for an existing service principal, " +
			"which will be granted access to the new vault.",
	}

	p.Properties["clientId"] = service.Parameter{
		Type: "string",
		Description: "Client ID (username) for an existing service principal," +
			"which will be granted access to the new vault.",
	}

	p.Properties["clientSecret"] = service.Parameter{
		Type: "string",
		Description: "Client secret (password) for an existing service " +
			"principal, which will be granted access to the new vault.",
	}

	p.Required = []string{"objectId", "clientId", "clientSecret"}
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
