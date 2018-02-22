package eventhubs

import "github.com/Azure/open-service-broker-azure/pkg/service"

// ProvisioningParameters encapsulates non-sensitive Azure Event Hub
// provisioning options
type ProvisioningParameters struct{}

// SecureProvisioningParameters encapsulates sensitive Azure Event Hub
// provisioning options
type SecureProvisioningParameters struct{}

type eventHubInstanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	EventHubName      string `json:"eventHubName"`
	EventHubNamespace string `json:"eventHubNamespace"`
}

type eventHubSecureInstanceDetails struct {
	ConnectionString string `json:"connectionString"`
	PrimaryKey       string `json:"primaryKey"`
}

// UpdatingParameters encapsulates search-specific updating options
type UpdatingParameters struct {
}

// BindingParameters encapsulates Azure Event Hub specific binding options
type BindingParameters struct {
}

type eventHubBindingDetails struct {
}

type eventHubSecureBindingDetails struct {
}

// Credentials encapsulates Event Hub-specific coonection details and
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
) GetEmptySecureProvisioningParameters() service.SecureProvisioningParameters {
	return &SecureProvisioningParameters{}
}

func (
	s *serviceManager,
) GetEmptyUpdatingParameters() service.UpdatingParameters {
	return &UpdatingParameters{}
}

func (
	s *serviceManager,
) GetEmptyInstanceDetails() service.InstanceDetails {
	return &eventHubInstanceDetails{}
}

func (
	s *serviceManager,
) GetEmptySecureInstanceDetails() service.SecureInstanceDetails {
	return &eventHubSecureInstanceDetails{}
}

func (s *serviceManager) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (s *serviceManager) GetEmptyBindingDetails() service.BindingDetails {
	return &eventHubBindingDetails{}
}

func (
	s *serviceManager,
) GetEmptySecureBindingDetails() service.SecureBindingDetails {
	return &eventHubSecureBindingDetails{}
}
