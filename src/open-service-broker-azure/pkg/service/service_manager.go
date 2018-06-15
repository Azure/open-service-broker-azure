package service

// ServiceManager is an interface to be implemented by module components
// responsible for managing the lifecycle of services and plans thereof
type ServiceManager interface { // nolint: golint
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
	// Bind synchronously binds to a service
	Bind(
		Instance,
		BindingParameters,
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
