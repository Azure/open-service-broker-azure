package fake

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

// ProvisioningValidationFunction describes a function used to provide pluggable
// provisioning validation behavior to the fake implementation of the
// service.Module interface
type ProvisioningValidationFunction func(service.ProvisioningParameters) error

// UpdatingValidationFunction describes a function used to provide pluggable
// updating validation behavior to the fake implementation of the
// service.Module interface
type UpdatingValidationFunction func(service.UpdatingParameters) error

// BindingValidationFunction describes a function used to provide pluggable
// binding validation behavior to the fake implementation of the service.Module
// interface
type BindingValidationFunction func(service.BindingParameters) error

// BindFunction describes a function used to provide pluggable binding behavior
// to the fake implementation of the service.Module interface
type BindFunction func(
	service.StandardProvisioningContext,
	service.ProvisioningContext,
	service.BindingParameters,
) (service.BindingContext, service.Credentials, error)

// UnbindFunction describes a function used to provide pluggable unbinding
// behavior to the fake implementation of the service.Module interface
type UnbindFunction func(
	service.StandardProvisioningContext,
	service.ProvisioningContext,
	service.BindingContext,
) error

// Module is a fake implementation of the service.Module interface used to
// facilittate testing.
type Module struct {
	ServiceManager *ServiceManager
}

// ServiceManager is a fake implementation of the service.ServiceManager
// interface used to facilitate testing.
type ServiceManager struct {
	ProvisioningValidationBehavior ProvisioningValidationFunction
	UpdatingValidationBehavior     UpdatingValidationFunction
	BindingValidationBehavior      BindingValidationFunction
	BindBehavior                   BindFunction
	UnbindBehavior                 UnbindFunction
}

// New returns a new instance of a type that fulfills the service.Module
// and provides an example of how such a module is implemented
func New() (*Module, error) {
	return &Module{
		ServiceManager: &ServiceManager{
			ProvisioningValidationBehavior: defaultProvisioningValidationBehavior,
			UpdatingValidationBehavior:     defaultUpdatingValidationBehavior,
			BindingValidationBehavior:      defaultBindingValidationBehavior,
			BindBehavior:                   defaultBindBehavior,
			UnbindBehavior:                 defaultUnbindBehavior,
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

// ValidateProvisioningParameters validates the provided provisioningParameters
// and returns an error if there is any problem
func (s *ServiceManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	return s.ProvisioningValidationBehavior(provisioningParameters)
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
	_ service.Plan,
) (service.ProvisioningContext, error) {
	return instance.ProvisioningContext, nil
}

// ValidateUpdatingParameters validates the provided updatingParameters
// and returns an error if there is any problem
func (s *ServiceManager) ValidateUpdatingParameters(
	updatingParameters service.UpdatingParameters,
) error {
	return s.UpdatingValidationBehavior(updatingParameters)
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
	_ service.Plan,
) (service.ProvisioningContext, error) {
	return instance.ProvisioningContext, nil
}

// ValidateBindingParameters validates the provided bindingParameters and
// returns an error if there is any problem
func (s *ServiceManager) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	return s.BindingValidationBehavior(bindingParameters)
}

// Bind synchronously binds to a service
func (s *ServiceManager) Bind(
	instance service.Instance,
	bindingParameters service.BindingParameters,
) (service.BindingContext, service.Credentials, error) {
	return s.BindBehavior(
		instance.StandardProvisioningContext,
		instance.ProvisioningContext,
		bindingParameters,
	)
}

// Unbind synchronously unbinds from a service
func (s *ServiceManager) Unbind(
	instance service.Instance,
	bindingContext service.BindingContext,
) error {
	return s.UnbindBehavior(
		instance.StandardProvisioningContext,
		instance.ProvisioningContext,
		bindingContext,
	)
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
	_ service.Plan,
) (service.ProvisioningContext, error) {
	return instance.ProvisioningContext, nil
}

func defaultProvisioningValidationBehavior(
	service.ProvisioningParameters,
) error {
	return nil
}

func defaultUpdatingValidationBehavior(
	service.UpdatingParameters,
) error {
	return nil
}

func defaultBindingValidationBehavior(service.BindingParameters) error {
	return nil
}

func defaultBindBehavior(
	_ service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	bindingParameters service.BindingParameters,
) (service.BindingContext, service.Credentials, error) {
	return provisioningContext, &Credentials{}, nil
}

func defaultUnbindBehavior(
	_ service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	bindingContext service.BindingContext,
) error {
	return nil
}
