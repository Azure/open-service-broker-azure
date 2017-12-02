package servicebus

import "github.com/Azure/open-service-broker-azure/pkg/service"

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
) GetEmptyProvisioningContext() service.ProvisioningContext {
	return &serviceBusProvisioningContext{}
}

func (s *serviceManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (s *serviceManager) GetEmptyBindingContext() service.BindingContext {
	return &serviceBusBindingContext{}
}

func (s *serviceManager) GetEmptyCredentials() service.Credentials {
	return &Credentials{}
}
