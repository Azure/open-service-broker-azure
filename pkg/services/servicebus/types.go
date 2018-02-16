package servicebus

import "github.com/Azure/open-service-broker-azure/pkg/service"

// ProvisioningParameters encapsulates Service Bus specific provisioning options
type ProvisioningParameters struct{}

type serviceBusInstanceDetails struct {
	ARMDeploymentName       string `json:"armDeployment"`
	ServiceBusNamespaceName string `json:"serviceBusNamespaceName"`
}

type serviceBusSecureInstanceDetails struct {
	ConnectionString string `json:"connectionString"`
	PrimaryKey       string `json:"primaryKey"`
}

// UpdatingParameters encapsulates servicebus-specific updating options
type UpdatingParameters struct {
}

// BindingParameters encapsulates Service Bus specific binding options
type BindingParameters struct {
}

type serviceBusBindingDetails struct {
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
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &serviceBusInstanceDetails{}
}

func (
	s *serviceManager,
) GetEmptySecureInstanceDetails() service.SecureInstanceDetails {
	return &serviceBusSecureInstanceDetails{}
}

func (s *serviceManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (s *serviceManager) GetEmptyBindingDetails() service.BindingDetails {
	return &serviceBusBindingDetails{}
}
