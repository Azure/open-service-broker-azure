package eventhub

import "github.com/Azure/azure-service-broker/pkg/service"

// ProvisioningParameters encapsulates Azure Event Hub provisioning options
type ProvisioningParameters struct{}

type eventHubProvisioningContext struct {
	ARMDeploymentName string `json:"armDeployment"`
	EventHubName      string `json:"eventHubName"`
	EventHubNamespace string `json:"eventHubNamespace"`
	PrimaryKey        string `json:"primaryKey"`
	ConnectionString  string `json:"connectionString"`
}

// UpdatingParameters encapsulates search-specific updating options
type UpdatingParameters struct {
}

// BindingParameters encapsulates Azure Event Hub specific binding options
type BindingParameters struct {
}

type eventHubBindingContext struct {
}

// Credentials encapsulates Event Hub-specific coonection details and
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
	return &eventHubProvisioningContext{}
}

func (m *module) GetEmptyBindingParameters() service.BindingParameters {
	return &BindingParameters{}
}

func (m *module) GetEmptyBindingContext() service.BindingContext {
	return &eventHubBindingContext{}
}

func (m *module) GetEmptyCredentials() service.Credentials {
	return &Credentials{}
}
