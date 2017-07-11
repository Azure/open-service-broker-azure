package service

import (
	"context"
	"fmt"
)

// ProvisioningStepFunction is the signature for functions that implement a
// provisioning step
type ProvisioningStepFunction func(
	ctx context.Context,
	provisioningResult interface{},
	params interface{},
) (interface{}, error)

// ProvisioningStep is an interface to be implemented by types that represent
// a single step in a chain of steps that defines a provisioning process
type ProvisioningStep interface {
	GetName() string
	Execute(
		ctx context.Context,
		provisioningResult,
		params interface{},
	) (interface{}, error)
}

type provisioningStep struct {
	name string
	fn   ProvisioningStepFunction
}

// Provisioner is an interface to be implemented by types that model a declared
// chain of tasks used to asynchronously provision a service
type Provisioner interface {
	GetFirstStepName() (string, bool)
	GetStep(name string) (ProvisioningStep, bool)
	GetNextStepName(name string) (string, bool)
}

type provisioner struct {
	firstStepName string
	steps         map[string]ProvisioningStep
	nextSteps     map[string]string
}

// NewProvisioningStep returns a new ProvisioningStep
func NewProvisioningStep(
	name string,
	fn ProvisioningStepFunction,
) ProvisioningStep {
	return &provisioningStep{
		name: name,
		fn:   fn,
	}
}

// GetName returns a provisioning step's name
func (p *provisioningStep) GetName() string {
	return p.name
}

// Execute executes a step
func (p *provisioningStep) Execute(
	ctx context.Context,
	provisioningResult interface{},
	params interface{},
) (interface{}, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	return p.fn(ctx, provisioningResult, params)
}

// NewProvisioner returns a new provisioner
func NewProvisioner(steps ...ProvisioningStep) (Provisioner, error) {
	p := &provisioner{
		steps:     make(map[string]ProvisioningStep),
		nextSteps: make(map[string]string),
	}
	if len(steps) > 0 {
		p.firstStepName = steps[0].GetName()
		var lastStep ProvisioningStep
		for _, step := range steps {
			_, ok := p.steps[step.GetName()]
			if ok {
				// This means a duplicate step name has been detected. This is a serious
				// problem.
				return nil, fmt.Errorf(
					`duplicate step name "%s" detected`,
					step.GetName(),
				)
			}
			p.steps[step.GetName()] = step
			if lastStep != nil {
				p.nextSteps[lastStep.GetName()] = step.GetName()
			}
			lastStep = step
		}
	}
	return p, nil
}

// GetFirstStepName retrieves the name of the first step in the chain
func (p *provisioner) GetFirstStepName() (string, bool) {
	return p.firstStepName, (p.firstStepName != "")
}

// GetStep retrieves a step by name
func (p *provisioner) GetStep(name string) (ProvisioningStep, bool) {
	step, ok := p.steps[name]
	return step, ok
}

// GetNextStepName, given the name of one step, returns the name of the next
// step and a boolean indicating whether a next step actually exists
func (p *provisioner) GetNextStepName(name string) (string, bool) {
	nextStepName, ok := p.nextSteps[name]
	return nextStepName, ok
}
