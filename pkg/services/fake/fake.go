package fake

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

// ProvisioningValidationFunction describes a function used to provide pluggable
// provisioning validation behavior to the fake implementation of the
// service.Module interface
type ProvisioningValidationFunction func(
	service.Plan,
	service.ProvisioningParameters,
	service.SecureProvisioningParameters,
) error

// BindFunction describes a function used to provide pluggable binding behavior
// to the fake implementation of the service.Module interface
type BindFunction func(
	service.Instance,
	service.BindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error)

// UnbindFunction describes a function used to provide pluggable unbinding
// behavior to the fake implementation of the service.Module interface
type UnbindFunction func(
	service.Instance,
	service.Binding,
) error

// Module is a fake implementation of the service.Module interface used to
// facilittate testing.
type Module struct {
	ServiceManager *ServiceManager
}

// ServiceManager is a fake implementation of the service.ServiceManager
// interface used to facilitate testing.
type ServiceManager struct {
	BindBehavior   BindFunction
	UnbindBehavior UnbindFunction
}

// New returns a new instance of a type that fulfills the service.Module
// and provides an example of how such a module is implemented
func New() (*Module, error) {
	return &Module{
		ServiceManager: &ServiceManager{
			BindBehavior:   defaultBindBehavior,
			UnbindBehavior: defaultUnbindBehavior,
		},
	}, nil
}

// GetName returns this module's name
func (m *Module) GetName() string {
	return "fake"
}

// GetStability returns this module's relative stability
func (m *Module) GetStability() service.Stability {
	return service.StabilityStable
}

// GetProvisioner returns a provisioner that defines the steps a module must
// execute asynchronously to provision a service
func (s *ServiceManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("run", s.provision),
	)
}

func (s *ServiceManager) provision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	return instance.Details, instance.SecureDetails, nil
}

// GetUpdater returns a updater that defines the steps a module must
// execute asynchronously to update a service
func (s *ServiceManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("run", s.update),
	)
}

func (s *ServiceManager) update(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	return instance.Details, instance.SecureDetails, nil
}

// Bind synchronously binds to a service
func (s *ServiceManager) Bind(
	instance service.Instance,
	bindingParameters service.BindingParameters,
	_ service.SecureBindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	return s.BindBehavior(instance, bindingParameters)
}

// GetCredentials returns service-specific credentials populated from instance
// and binding details
func (s *ServiceManager) GetCredentials(
	service.Instance,
	service.Binding,
) (service.Credentials, error) {
	return nil, nil
}

// Unbind synchronously unbinds from a service
func (s *ServiceManager) Unbind(
	instance service.Instance,
	binding service.Binding,
) error {
	return s.UnbindBehavior(instance, binding)
}

// GetDeprovisioner returns a deprovisioner that defines the steps a module
// must execute asynchronously to deprovision a service
func (s *ServiceManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("run", s.deprovision),
	)
}

func (s *ServiceManager) deprovision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	return instance.Details, instance.SecureDetails, nil
}

func defaultBindBehavior(
	service.Instance,
	service.BindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	return nil, nil, nil
}

func defaultUnbindBehavior(
	service.Instance,
	service.Binding,
) error {
	return nil
}
