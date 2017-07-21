package service

// Module is an interface to be implemented by the broker's service modules
type Module interface {
	// GetName returns a module's name
	GetName() string
	// GetCatalog returns a Catalog of service/plans offered by a module
	GetCatalog() (Catalog, error)
	// GetEmptyProvisioningParameters returns an empty instance of module-specific
	// provisioningParameters
	GetEmptyProvisioningParameters() interface{}
	// ValidateProvisioningParameters validates the provided
	// provisioningParameters and returns an error if there is any problem
	ValidateProvisioningParameters(params interface{}) error
	// GetProvisioner returns a provisioner that defines the steps a module must
	// execute asynchronously to provision a service
	GetProvisioner() (Provisioner, error)
	// GetEmptyProvisioningContext returns an empty instance of a module-specific
	// provisioningContext
	GetEmptyProvisioningContext() interface{}
	// GetEmptyBindingParameters returns an empty instance of module-specific
	// bindingParameters
	GetEmptyBindingParameters() interface{}
	// ValidateBindingParameters validates the provided bindingParameters and
	// returns an error if there is any problem
	ValidateBindingParameters(params interface{}) error
	// Bind synchronously binds to a service
	Bind(provisioningContext, params interface{}) (interface{}, error)
	// GetEmptyBindingContext returns an empty instance of a module-specific
	// bindingContext
	GetEmptyBindingContext() interface{}
	// Unbind synchronously unbinds from a service
	Unbind(provisioningContext, bindingContext interface{}) error
	// GetDeprovisioner returns a deprovisioner that defines the steps a module
	// must execute asynchronously to deprovision a service
	GetDeprovisioner() (Deprovisioner, error)
}
