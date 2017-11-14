package keyvault

import "github.com/Azure/azure-service-broker/pkg/service"

// ProvisioningParameters encapsulates keyvault-specific provisioning options
type ProvisioningParameters struct {
	Location      string            `json:"location"`
	ResourceGroup string            `json:"resourceGroup"`
	Tags          map[string]string `json:"tags"`
	ObjectID      string            `json:"objectId"`
	ClientID      string            `json:"clientId"`
	ClientSecret  string            `json:"clientSecret"`
}

type keyvaultProvisioningContext struct {
	ResourceGroupName string `json:"resourceGroup"`
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

type keyvaultBindingContext struct {
}

// Credentials encapsulates Key Vault-specific coonection details and
// credentials.
type Credentials struct {
	VaultURI     string `json:"vaultUri"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

func (
	m *module,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ProvisioningParameters{}
}

// SetResourceGroup sets the name of the resource group into which service
// instances will be deployed
func (p *ProvisioningParameters) SetResourceGroup(resourceGroup string) {
	p.ResourceGroup = resourceGroup
}

func (
	m *module,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (m *module) GetEmptyProvisioningContext() service.ProvisioningContext {
	return &keyvaultProvisioningContext{}
}

func (m *module) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (m *module) GetEmptyBindingContext() service.BindingContext {
	return &keyvaultBindingContext{}
}

func (m *module) GetEmptyCredentials() service.Credentials {
	return &Credentials{}
}
