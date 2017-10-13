package keyvault

import "github.com/Azure/azure-service-broker/pkg/service"

// ProvisioningParameters encapsulates keyvault-specific provisioning options
type ProvisioningParameters struct {
	Location      string            `json:"location"`
	ResourceGroup string            `json:"resourceGroup"`
	Tags          map[string]string `json:"tags"`
	ObjectID      string            `json:"objectid"`
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

func (k *keyvaultProvisioningContext) GetResourceGroupName() string {
	return k.ResourceGroupName
}

func (
	m *module,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ProvisioningParameters{}
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
