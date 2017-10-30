package client

// ProvisioningParameters is a specialized map[string]interface{} that
// implements service.ProvisioningParameters. Unlike a module-specific
// implementation of service.ProvisioningParameters used server-side, this
// implementation has the flexibility to encapsultate parameters for ANY
// service/plan.
type ProvisioningParameters map[string]interface{}

// SetResourceGroup sets the resourceGroup to be included in these
// provisioningParameters. This isn't really used, but it exists to ensure
// ProvisioningParameters conforms to the service.ProvisioningParameters
// interface.
func (p ProvisioningParameters) SetResourceGroup(resourceGroup string) {
	p["resourceGroup"] = resourceGroup
}
