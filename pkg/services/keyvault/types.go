package keyvault

import "github.com/Azure/open-service-broker-azure/pkg/service"

// ProvisioningParameters encapsulates keyvault-specific provisioning options
type ProvisioningParameters struct {
	ObjectID     string `json:"objectId"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type keyvaultInstanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	KeyVaultName      string `json:"keyVaultName"`
	VaultURI          string `json:"vaultUri"`
	ClientID          string `json:"clientId"`
	ClientSecret      string `json:"clientSecret"`
}

// UpdatingParameters encapsulates keyvault-specific updating options
type UpdatingParameters struct {
}

// BindingParameters encapsulates keyvault-specific binding options
type BindingParameters struct {
}

type keyvaultBindingDetails struct {
}

// Credentials encapsulates Key Vault-specific coonection details and
// credentials.
type Credentials struct {
	VaultURI     string `json:"vaultUri"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

func (
	s *serviceManager,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ProvisioningParameters{}
}

func (
	s *serviceManager,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (
	s *serviceManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &keyvaultInstanceDetails{}
}

func (s *serviceManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (s *serviceManager) GetEmptyBindingDetails() service.BindingDetails {
	return &keyvaultBindingDetails{}
}
