package async

import (
	"context"
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/machinery"
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/storage"
)

const (
	operationProvisioning   = "provisioning"
	operationDeprovisioning = "deprovisioning"
)

// Engine is an interface for a broker-specifc framework for submitting and
// asynchronously completing provisioning and deprovisioning tasks.
// Implementations of the Enginer interface are abstractions over the
// complexities of the underlying async task execution library.
type Engine interface {
	// Provision kicks off asynchronous provisioning for the given instance
	Provision(*service.Instance) error
	// Deprovision kicks off asynchronous deprovisioning for the given instance
	Deprovision(*service.Instance) error
	// Start causes the async engine to begin executing queued tasks
	Start(context.Context) error
}

// engine is a machinery-based implementation of the Engine interface.
type engine struct {
	store           storage.Store
	machineryServer machinery.Server
	modules         map[string]service.Module
	provisioners    map[string]service.Provisioner
	deprovisioners  map[string]service.Deprovisioner
	// This allows tests to inject an alternative implementation of this function
	getWorker func(machinery.Server) machinery.Worker
}

// NewEngine returns a new machinery-based implementation of the Engine
// interface
func NewEngine(
	store storage.Store,
	machineryServer machinery.Server,
	modules []service.Module,
) (Engine, error) {
	e := &engine{
		store:           store,
		machineryServer: machineryServer,
		modules:         make(map[string]service.Module),
		provisioners:    make(map[string]service.Provisioner),
		deprovisioners:  make(map[string]service.Deprovisioner),
		getWorker:       getWorker,
	}
	machineryServer.RegisterTask("work", e.doWork)

	for _, module := range modules {
		catalog, err := module.GetCatalog()
		if err != nil {
			return nil, err
		}
		for _, svc := range catalog.GetServices() {
			existingModule, ok := e.modules[svc.GetID()]
			if ok {
				// This means we have more than one module claiming to provide services
				// with an ID in common. This is a SERIOUS problem.
				return nil, fmt.Errorf(
					"module %s and module %s BOTH provide a service with the id %s",
					existingModule.GetName(),
					module.GetName(),
					svc.GetID())
			}
			e.modules[svc.GetID()] = module
			provisioner, err := module.GetProvisioner()
			if err != nil {
				return nil, err
			}
			e.provisioners[svc.GetID()] = provisioner
			deprovisioner, err := module.GetDeprovisioner()
			if err != nil {
				return nil, err
			}
			e.deprovisioners[svc.GetID()] = deprovisioner
		}
	}

	return e, nil
}

func (e *engine) Start(ctx context.Context) error {
	errChan := make(chan error)
	defer close(errChan)
	go func() {
		worker := e.getWorker(e.machineryServer)
		errChan <- worker.Launch()
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}

func getWorker(machineryServer machinery.Server) machinery.Worker {
	return machineryServer.NewWorker("worker")
}
