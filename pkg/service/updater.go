package service

import (
	"context"
	"fmt"
)

// UpdatingStepFunction is the signature for functions that implement a
// updating step
type UpdatingStepFunction func(
	ctx context.Context,
	instanceID string,
	serviceID string,
	planID string,
	standardProvisioningContext StandardProvisioningContext,
	provisioningContext ProvisioningContext,
	params UpdatingParameters,
) (ProvisioningContext, error)

// UpdatingStep is an interface to be implemented by types that represent
// a single step in a chain of steps that defines a updating process
type UpdatingStep interface {
	GetName() string
	Execute(
		ctx context.Context,
		instanceID string,
		serviceID string,
		planID string,
		standardProvisioningContext StandardProvisioningContext,
		provisioningContext ProvisioningContext,
		params UpdatingParameters,
	) (ProvisioningContext, error)
}

type updatingStep struct {
	name string
	fn   UpdatingStepFunction
}

// Updater is an interface to be implemented by types that model a declared
// chain of tasks used to asynchronously update a service
type Updater interface {
	GetFirstStepName() (string, bool)
	GetStep(name string) (UpdatingStep, bool)
	GetNextStepName(name string) (string, bool)
}

type updater struct {
	firstStepName string
	steps         map[string]UpdatingStep
	nextSteps     map[string]string
}

// NewUpdatingStep returns a new UpdatingStep
func NewUpdatingStep(
	name string,
	fn UpdatingStepFunction,
) UpdatingStep {
	return &updatingStep{
		name: name,
		fn:   fn,
	}
}

// GetName returns a updating step's name
func (u *updatingStep) GetName() string {
	return u.name
}

// Execute executes a step
func (u *updatingStep) Execute(
	ctx context.Context,
	instanceID string,
	serviceID string,
	planID string,
	standardProvisioningContext StandardProvisioningContext,
	provisioningContext ProvisioningContext,
	params UpdatingParameters,
) (ProvisioningContext, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	return u.fn(
		ctx,
		instanceID,
		serviceID,
		planID,
		standardProvisioningContext,
		provisioningContext,
		params,
	)
}

// NewUpdater returns a new updater
func NewUpdater(steps ...UpdatingStep) (Updater, error) {
	u := &updater{
		steps:     make(map[string]UpdatingStep),
		nextSteps: make(map[string]string),
	}
	if len(steps) > 0 {
		u.firstStepName = steps[0].GetName()
		var lastStep UpdatingStep
		for _, step := range steps {
			_, ok := u.steps[step.GetName()]
			if ok {
				// This means a duplicate step name has been detected. This is a serious
				// problem.
				return nil, fmt.Errorf(
					`duplicate step name "%s" detected`,
					step.GetName(),
				)
			}
			u.steps[step.GetName()] = step
			if lastStep != nil {
				u.nextSteps[lastStep.GetName()] = step.GetName()
			}
			lastStep = step
		}
	}
	return u, nil
}

// GetFirstStepName retrieves the name of the first step in the chain
func (u *updater) GetFirstStepName() (string, bool) {
	return u.firstStepName, (u.firstStepName != "")
}

// GetStep retrieves a step by name
func (u *updater) GetStep(name string) (UpdatingStep, bool) {
	step, ok := u.steps[name]
	return step, ok
}

// GetNextStepName, given the name of one step, returns the name of the next
// step and a boolean indicating whether a next step actually exists
func (u *updater) GetNextStepName(name string) (string, bool) {
	nextStepName, ok := u.nextSteps[name]
	return nextStepName, ok
}
