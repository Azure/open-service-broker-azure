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

// BindingParameters encapsulates Service Bus specific binding options
type BindingParameters struct {
}

type serviceBusBindingContext struct {
}

type serviceBusCredentials struct {
	ConnectionString string `json:"connectionString"`
	PrimaryKey       string `json:"primaryKey"`
}

func (r *serviceBusProvisioningContext) GetResourceGroupName() string {
	return r.ResourceGroupName
}

func (
	m *module,
) GetEmptyProvisioningParameters() service.ProvisioningParameters {
	return &ProvisioningParameters{}
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
	return &serviceBusCredentials{}
}
