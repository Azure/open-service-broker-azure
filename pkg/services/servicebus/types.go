package servicebus

import "github.com/Azure/azure-service-broker/pkg/service"

// ProvisioningParameters encapsulates Service Bus specific provisioning options
type ProvisioningParameters struct{}

type serviceBusProvisioningContext struct {
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
