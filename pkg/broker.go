package main

import (
	"context"
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/api"
	"github.com/Azure/azure-service-broker/pkg/async"
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/storage"
	"github.com/go-redis/redis"
)

type errAsyncEngineStopped struct {
	err error
}

func (e *errAsyncEngineStopped) Error() string {
	return fmt.Sprintf("async engine stopped: %s", e.err)
}

type errAPIServerStopped struct {
	err error
}

func (e *errAPIServerStopped) Error() string {
	return fmt.Sprintf("api server stopped: %s", e.err)
}

type broker struct {
	store       storage.Store
	apiServer   api.Server
	asyncEngine async.Engine
	// Modules indexed by service
	modules map[string]service.Module
	// Provisioners indexed by service
	provisioners map[string]service.Provisioner
	// Deprovisioners indexed by service
	deprovisioners map[string]service.Deprovisioner
}

func newBroker(
	redisClient *redis.Client,
	modules []service.Module,
) (*broker, error) {
	b := &broker{
		store:          storage.NewStore(redisClient),
		asyncEngine:    async.NewEngine(redisClient),
		modules:        make(map[string]service.Module),
		provisioners:   make(map[string]service.Provisioner),
		deprovisioners: make(map[string]service.Deprovisioner),
	}

	for _, module := range modules {
		catalog, err := module.GetCatalog()
		if err != nil {
			return nil, err
		}
		for _, svc := range catalog.GetServices() {
			existingModule, ok := b.modules[svc.GetID()]
			if ok {
				// This means we have more than one module claiming to provide services
				// with an ID in common. This is a SERIOUS problem.
				return nil, fmt.Errorf(
					"module %s and module %s BOTH provide a service with the id %s",
					existingModule.GetName(),
					module.GetName(),
					svc.GetID())
			}
			b.modules[svc.GetID()] = module
			provisioner, err := module.GetProvisioner()
			if err != nil {
				return nil, err
			}
			b.provisioners[svc.GetID()] = provisioner
			deprovisioner, err := module.GetDeprovisioner()
			if err != nil {
				return nil, err
			}
			b.deprovisioners[svc.GetID()] = deprovisioner
		}
	}

	b.asyncEngine.RegisterJob("provisionStep", b.doProvisionStep)
	b.asyncEngine.RegisterJob("deprovisionStep", b.doDeprovisionStep)

	var err error
	b.apiServer, err = api.NewServer(
		8080,
		storage.NewStore(redisClient),
		b.asyncEngine,
		b.modules,
		b.provisioners,
		b.deprovisioners,
	)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (b *broker) start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errChan := make(chan error)
	// Start async engine
	go func() {
		err := b.asyncEngine.Start(ctx)
		aes := &errAsyncEngineStopped{err: err}
		select {
		case errChan <- aes:
		case <-ctx.Done():
		}
	}()
	// Start api server
	go func() {
		err := b.apiServer.Start(ctx)
		ss := &errAPIServerStopped{err: err}
		select {
		case errChan <- ss:
		case <-ctx.Done():
		}
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}
