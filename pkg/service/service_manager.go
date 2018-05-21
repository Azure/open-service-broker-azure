package service

// ServiceManager is an interface to be implemented by module components
// responsible for managing the lifecycle of services and plans thereof
type ServiceManager interface { // nolint: golint
	// SplitProvisioningParameters splits a map of provisioning parameters into
	// two separate maps, with one containing non-sensitive provisioning
	// parameters and the other containing sensitive provisioning parameters.
	SplitProvisioningParameters(
		map[string]interface{},
	) (ProvisioningParameters, SecureProvisioningParameters, error)
	// GetProvisioner returns a provisioner that defines the steps a module must
	// execute asynchronously to provision a service.
	GetProvisioner(Plan) (Provisioner, error)
	// ValidateUpdatingParameters validates the provided
	// updating parameters against against current instance state
	// and returns an error if there is any problem
	ValidateUpdatingParameters(Instance) error
	// GetUpdater returns a updater that defines the steps a module must
	// execute asynchronously to update a service.
	GetUpdater(Plan) (Updater, error)
	// SplitBindingParameters splits a map of binding parameters into two separate
	// maps, with one containing non-sensitive binding parameters and the other
	// containing sensitive binding parameters.
	SplitBindingParameters(
		CombinedBindingParameters,
	) (BindingParameters, SecureBindingParameters, error)
	// Bind synchronously binds to a service
	Bind(
		Instance,
		BindingParameters,
		SecureBindingParameters,
	) (BindingDetails, SecureBindingDetails, error)
	// GetCredentials returns service-specific credentials populated from instance
	// and binding details
	GetCredentials(Instance, Binding) (Credentials, error)
	// Unbind synchronously unbinds from a service
	Unbind(Instance, Binding) error
	// GetDeprovisioner returns a deprovisioner that defines the steps a module
	// must execute asynchronously to deprovision a service
	GetDeprovisioner(Plan) (Deprovisioner, error)
}
