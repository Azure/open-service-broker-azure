package fake

import (
	"context"

	"github.com/Azure/azure-service-broker/pkg/service"
)

// Engine is a fake implementation of async.Engine used for testing
type Engine struct {
	RunBehavior func() error
}

// NewEngine returns a new, fake implementation of async.Engine used for testing
func NewEngine() *Engine {
	return &Engine{
		RunBehavior: defaultRunBehavior,
	}
}

// Provision kicks off asynchronous provisioning for the given instance
func (e *Engine) Provision(instance *service.Instance) error {
	return nil
}

// Deprovision kicks off asynchronous deprovisioning for the given instance
func (e *Engine) Deprovision(instance *service.Instance) error {
	return nil
}

// Start causes the async engine to begin executing queued tasks
func (e *Engine) Start(context.Context) error {
	return e.RunBehavior()
}

func defaultRunBehavior() error {
	select {}
}
