package service

import (
	"context"
	"fmt"
)

// DeprovisioningStepFunction is the signature for functions that implement a
// deprovisioning step
type DeprovisioningStepFunction func(
	ctx context.Context,
	instance Instance,
	plan Plan,
	refInstance Instance,
) (InstanceDetails, error)

// DeprovisioningStep is an interface to be implemented by types that represent
// a single step in a chain of steps that defines a deprovisioning process
type DeprovisioningStep interface {
	GetName() string
	Execute(
		ctx context.Context,
		instance Instance,
		plan Plan,
		refInstance Instance,
	) (InstanceDetails, error)
}

type deprovisioningStep struct {
	name string
	fn   DeprovisioningStepFunction
}

// Deprovisioner is an interface to be implemented by types that model a
// declared chain of tasks used to asynchronously deprovision a service
type Deprovisioner interface {
	GetFirstStepName() (string, bool)
	GetStep(name string) (DeprovisioningStep, bool)
	GetNextStepName(name string) (string, bool)
}

type deprovisioner struct {
	firstStepName string
	steps         map[string]DeprovisioningStep
	nextSteps     map[string]string
}

// NewDeprovisioningStep returns a new DeprovisioningStep
func NewDeprovisioningStep(
	name string,
	fn DeprovisioningStepFunction,
) DeprovisioningStep {
	return &deprovisioningStep{
		name: name,
		fn:   fn,
	}
}

// GetName returns a deprovisioning step's name
func (d *deprovisioningStep) GetName() string {
	return d.name
}

// Execute executes a step
func (d *deprovisioningStep) Execute(
	ctx context.Context,
	instance Instance,
	plan Plan,
	refInstance Instance,
) (InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	return d.fn(
		ctx,
		instance,
		plan,
		refInstance,
	)
}

// NewDeprovisioner returns a new deprovisioner
func NewDeprovisioner(steps ...DeprovisioningStep) (Deprovisioner, error) {
	d := &deprovisioner{
		steps:     make(map[string]DeprovisioningStep),
		nextSteps: make(map[string]string),
	}
	if len(steps) > 0 {
		d.firstStepName = steps[0].GetName()
		var lastStep DeprovisioningStep
		for _, step := range steps {
			_, ok := d.steps[step.GetName()]
			if ok {
				// This means a duplicate step name has been detected. This is a serious
				// problem.
				return nil, fmt.Errorf(
					`duplicate step name "%s" detected`,
					step.GetName(),
				)
			}
			d.steps[step.GetName()] = step
			if lastStep != nil {
				d.nextSteps[lastStep.GetName()] = step.GetName()
			}
			lastStep = step
		}
	}
	return d, nil
}

// GetFirstStepName retrieves the name of the first step in the chain
func (d *deprovisioner) GetFirstStepName() (string, bool) {
	return d.firstStepName, (d.firstStepName != "")
}

// GetStep retrieves a step by name
func (d *deprovisioner) GetStep(name string) (DeprovisioningStep, bool) {
	step, ok := d.steps[name]
	return step, ok
}

// GetNextStepName, given the name of one step, returns the name of the next
// step and a boolean indicating whether a next step actually exists
func (d *deprovisioner) GetNextStepName(name string) (string, bool) {
	nextStepName, ok := d.nextSteps[name]
	return nextStepName, ok
}
