package service

// Module is an interface to be implemented by the broker's service modules
type Module interface {
	// GetName returns a module's name
	GetName() string
	// GetStability returns a module's relative level of stability
	GetStability() Stability
	// GetCatalog returns a Catalog of service/plans offered by a module
	GetCatalog() (Catalog, error)
	// GetEmptyProvisioningParameters returns an empty instance of module-specific
	// provisioningParameters
	GetEmptyProvisioningParameters() ProvisioningParameters
	// ValidateProvisioningParameters validates the provided
	// provisioningParameters and returns an error if there is any problem
	ValidateProvisioningParameters(ProvisioningParameters) error
	// GetProvisioner returns a provisioner that defines the steps a module must
	// execute asynchronously to provision a service.
	//
	// The two input parameters represent the service ID and plan ID
	// (respectively). Using these parameters, module implementations can
	// choose to return different Provisioner implementations for different
	// services/plans
	GetProvisioner(serviceID, planID string) (Provisioner, error)
	// GetEmptyProvisioningContext returns an empty instance of a module-specific
	// ProvisioningContext
	GetEmptyProvisioningContext() ProvisioningContext
	// GetEmptyUpdatingParameters returns an empty instance of module-specific
	// updatingParameters
	GetEmptyUpdatingParameters() UpdatingParameters
	// ValidateUpdatingParameters validates the provided
	// updatingParameters and returns an error if there is any problem
	ValidateUpdatingParameters(UpdatingParameters) error
	// GetUpdater returns a updater that defines the steps a module must
	// execute asynchronously to update a service.
	//
	// The two input parameters represent the service ID and plan ID
	// (respectively). Using these parameters, module implementations can
	// choose to return different Updater implementations for different
	// services/plans
	GetUpdater(serviceID, planID string) (Updater, error)
	// GetEmptyBindingParameters returns an empty instance of module-specific
	// bindingParameters
	GetEmptyBindingParameters() BindingParameters
	// ValidateBindingParameters validates the provided bindingParameters and
	// returns an error if there is any problem
	ValidateBindingParameters(BindingParameters) error
	// Bind synchronously binds to a service
	Bind(
		ProvisioningContext,
		BindingParameters,
	) (BindingContext, Credentials, error)
	// GetEmptyBindingContext returns an empty instance of a module-specific
	// bindingContext
	GetEmptyBindingContext() BindingContext
	// GetEmptyCredentials returns an empty instance of module-specific
	// credentials
	GetEmptyCredentials() Credentials
	// Unbind synchronously unbinds from a service
	Unbind(ProvisioningContext, BindingContext) error
	// GetDeprovisioner returns a deprovisioner that defines the steps a module
	// must execute asynchronously to deprovision a service
	//
	// The two input parameters represent the service ID and plan ID
	// (respectively). Using these parameters, module implementations can
	// choose to return different Deprovisioner implementations for different
	// services/plans
	GetDeprovisioner(serviceID, planID string) (Deprovisioner, error)
}
