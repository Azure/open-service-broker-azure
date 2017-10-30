package servicebus

import "github.com/Azure/azure-service-broker/pkg/service"

// ProvisioningParameters encapsulates Service Bus specific provisioning options
type ProvisioningParameters struct {
	Location      string            `json:"location"`
	ResourceGroup string            `json:"resourceGroup"`
	Tags          map[string]string `json:"tags"`
}

type serviceBusProvisioningContext struct {
	ResourceGroupName       string `json:"resourceGroup"`
	ARMDeploymentName       string `json:"armDeployment"`
	ServiceBusNamespaceName string `json:"serviceBusNamespaceName"`
	ConnectionString        string `json:"connectionString"`
	PrimaryKey              string `json:"primaryKey"`
}

// UpdatingParameters encapsulates servicebus-specific updating options
type UpdatingParameters struct {
}

// BindingParameters encapsulates Service Bus specific binding options
type BindingParameters struct {
}

type serviceBusBindingContext struct {
}

// Credentials encapsulates Service Bus-specific coonection details and
// credentials.
type Credentials struct {
	ConnectionString string `json:"connectionString"`
	PrimaryKey       string `json:"primaryKey"`
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
	return &serviceBusProvisioningContext{}
}

func (m *module) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (m *module) GetEmptyBindingContext() service.BindingContext {
	return &serviceBusBindingContext{}
}

func (m *module) GetEmptyCredentials() service.Credentials {
	return &Credentials{}
}
