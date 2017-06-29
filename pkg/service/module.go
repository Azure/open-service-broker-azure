package service

// Module is an interface to be implemented by the broker's service modules
type Module interface {
	// GetName returns a module's name
	GetName() string
	// GetCatalog returns a Catalog of service/plans offered by a module
	GetCatalog() (Catalog, error)
	// GetEmptyProvisioningParameters returns an empty instance of module-specific
	// provisioning parameters
	GetEmptyProvisioningParameters() interface{}
	// ValidateProvisioningParameters validates the provided provisioning
	// parameters and returns an error if there is any problem
	ValidateProvisioningParameters(params interface{}) error
	// GetProvisioner returns a provisioner that defines the steps a module must
	// execute asynchronously to provision a service
	GetProvisioner() (Provisioner, error)
	// GetEmptyProvisioningResult returns an empty instance of a module-specific
	// provisioning result
	GetEmptyProvisioningResult() interface{}
	// GetEmptyBindingParameters returns an empty instance of module-specific
	// binding parameters
	GetEmptyBindingParameters() interface{}
	// ValidateBindingParameters validates the provided binding parameters and
	// returns an error if there is any problem
	ValidateBindingParameters(params interface{}) error
	// GetEmptyBindingResult returns an empty instance of the type used to
	// encapsulate and store module-specific details of a completed binding
	// process
	Bind(provisioningResult, params interface{}) (interface{}, error)
	// GetEmptyBindingResult returns an empty instance of a module-specific
	// binding result
	GetEmptyBindingResult() interface{}
	// Unbind synchronously unbinds from a service
	Unbind(provisioningResult, bindingResult interface{}) error
	// GetDeprovisioner returns a deprovisioner that defines the steps a module
	// must execute asynchronously to deprovision a service
	GetDeprovisioner() (Deprovisioner, error)
}
