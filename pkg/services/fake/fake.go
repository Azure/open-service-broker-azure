package fake

import (
	"context"

	"github.com/Azure/azure-service-broker/pkg/service"
)

// ProvisioningValidationFunction describes a function used to provide pluggable
// provisioning validation behavior to the fake implementation of the
// service.Module interface
type ProvisioningValidationFunction func(service.ProvisioningParameters) error

// BindingValidationFunction describes a function used to provide pluggable
// binding validation behavior to the fake implementation of the service.Module
// interface
type BindingValidationFunction func(service.BindingParameters) error

// BindFunction describes a function used to provide pluggable binding behavior
// to the fake implementation of the service.Module interface
type BindFunction func(
	provisioningContext service.ProvisioningContext,
	bindingParameters service.BindingParameters,
) (service.BindingContext, service.Credentials, error)

// UnbindFunction describes a function used to provide pluggable unbinding
// behavior to the fake implementation of the service.Module interface
type UnbindFunction func(
	provisioningContext service.ProvisioningContext,
	bindingContext service.BindingContext,
) error

// Module is a fake implementation of the service.Module interface used to
// facilittate testing.
type Module struct {
	ProvisioningValidationBehavior ProvisioningValidationFunction
	BindingValidationBehavior      BindingValidationFunction
	BindBehavior                   BindFunction
	UnbindBehavior                 UnbindFunction
}

// New returns a new instance of a type that fulfills the service.Module
// and provides an example of how such a module is implemented
func New() (*Module, error) {
	return &Module{
		ProvisioningValidationBehavior: defaultProvisioningValidationBehavior,
		BindingValidationBehavior:      defaultBindingValidationBehavior,
		BindBehavior:                   defaultBindBehavior,
		UnbindBehavior:                 defaultUnbindBehavior,
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
func (m *Module) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	return m.ProvisioningValidationBehavior(provisioningParameters)
}

// GetProvisioner returns a provisioner that defines the steps a module must
// execute asynchronously to provision a service
func (m *Module) GetProvisioner(string, string) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("run", m.provision),
	)
}

func (m *Module) provision(
	ctx context.Context, // nolint: unparam
	instanceID string, // nolint: unparam
	serviceID string, // nolint: unparam
	planID string, // nolint: unparam
	provisioningContext service.ProvisioningContext,
	provisioningParameters service.ProvisioningParameters, // nolint: unparam
) (service.ProvisioningContext, error) {
	return provisioningContext, nil
}

// ValidateBindingParameters validates the provided bindingParameters and
// returns an error if there is any problem
func (m *Module) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	return m.BindingValidationBehavior(bindingParameters)
}

// Bind synchronously binds to a service
func (m *Module) Bind(
	provisioningContext service.ProvisioningContext,
	bindingParameters service.BindingParameters,
) (service.BindingContext, service.Credentials, error) {
	return m.BindBehavior(provisioningContext, bindingParameters)
}

// Unbind synchronously unbinds from a service
func (m *Module) Unbind(
	provisioningContext service.ProvisioningContext,
	bindingContext service.BindingContext,
) error {
	return m.UnbindBehavior(provisioningContext, bindingContext)
}

// GetDeprovisioner returns a deprovisioner that defines the steps a module
// must execute asynchronously to deprovision a service
func (m *Module) GetDeprovisioner(
	string,
	string,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("run", m.deprovision),
	)
}

func (m *Module) deprovision(
	ctx context.Context, // nolint: unparam
	instanceID string, // nolint: unparam
	serviceID string, // nolint: unparam
	planID string, // nolint: unparam
	provisioningContext service.ProvisioningContext,
) (service.ProvisioningContext, error) {
	return provisioningContext, nil
}

func defaultProvisioningValidationBehavior(
	service.ProvisioningParameters,
) error {
	return nil
}

func defaultBindingValidationBehavior(service.BindingParameters) error {
	return nil
}

func defaultBindBehavior(
	provisioningContext service.ProvisioningContext,
	bindingParameters service.BindingParameters,
) (service.BindingContext, service.Credentials, error) {
	return provisioningContext, &Credentials{}, nil
}

func defaultUnbindBehavior(
	provisioningContext service.ProvisioningContext,
	bindingContext service.BindingContext,
) error {
	return nil
}
