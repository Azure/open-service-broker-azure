package service

// ServiceManager is an interface to be implemented by module components
// responsible for managing the lifecycle of services and plans thereof
type ServiceManager interface { // nolint: golint
	// GetEmptyInstanceDetails returns an "empty" service-specific object that
	// can be populated with data during unmarshaling of JSON to an Instance
	GetEmptyInstanceDetails() InstanceDetails
	// GetProvisioner returns a provisioner that defines the steps a module must
	// execute asynchronously to provision a service.
	GetProvisioner(Plan) (Provisioner, error)
	// ValidateUpdatingParameters validates the provided
	// updating parameters against current instance state
	// and returns an error if there is any problem
	ValidateUpdatingParameters(Instance) error
	// GetUpdater returns a updater that defines the steps a module must
	// execute asynchronously to update a service.
	GetUpdater(Plan) (Updater, error)
	// GetEmptyBindingDetails returns an "empty" service-specific object that
	// can be populated with data during unmarshaling of JSON to a Binding
	GetEmptyBindingDetails() BindingDetails
	// Bind synchronously binds to a service
	Bind(Instance, BindingParameters) (BindingDetails, error)
	// GetCredentials returns service-specific credentials populated from instance
	// and binding details
	GetCredentials(Instance, Binding) (Credentials, error)
	// Unbind synchronously unbinds from a service
	Unbind(Instance, Binding) error
	// GetDeprovisioner returns a deprovisioner that defines the steps a module
	// must execute asynchronously to deprovision a service
	GetDeprovisioner(Plan) (Deprovisioner, error)
}
