package broker

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/api"
	"github.com/Azure/azure-service-broker/pkg/api/authenticator"
	"github.com/Azure/azure-service-broker/pkg/async"
	"github.com/Azure/azure-service-broker/pkg/crypto"
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/storage"
	log "github.com/Sirupsen/logrus"
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

// Broker is an interface to be implemented by components that implement full
// OSB functionality.
type Broker interface {
	// Start starts all broker components (e.g. API server and async execution
	// engine) and blocks until one of those components returns or fails.
	Start(context.Context) error
}

type broker struct {
	store       storage.Store
	apiServer   api.Server
	asyncEngine async.Engine
	codec       crypto.Codec
	// Modules indexed by service
	modules map[string]service.Module
}

// NewBroker returns a new Broker
func NewBroker(
	redisClient *redis.Client,
	codec crypto.Codec,
	authenticator authenticator.Authenticator,
	modules []service.Module,
	minStability service.Stability,
	defaultAzureLocation string,
	defaultAzureResourceGroup string,
) (Broker, error) {
	b := &broker{
		store:       storage.NewStore(redisClient),
		asyncEngine: async.NewEngine(redisClient),
		codec:       codec,
		modules:     make(map[string]service.Module),
	}

	for _, module := range modules {
		if module.GetStability() >= minStability {
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
						`module "%s" and module "%s" BOTH provide a service with the id "%s"`,
						existingModule.GetName(),
						module.GetName(),
						svc.GetID())
				}
				b.modules[svc.GetID()] = module
			}
		}
	}

	err := b.asyncEngine.RegisterJob("provisionStep", b.doProvisionStep)
	if err != nil {
		return nil, errors.New(
			"error registering async job for executing provisioning steps",
		)
	}
	err = b.asyncEngine.RegisterJob("updateStep", b.doUpdateStep)
	if err != nil {
		return nil, errors.New(
			"error registering async job for executing updating steps",
		)
	}
	err = b.asyncEngine.RegisterJob("deprovisionStep", b.doDeprovisionStep)
	if err != nil {
		return nil, errors.New(
			"error registering async job for executing deprovisioning steps",
		)
	}

	b.apiServer, err = api.NewServer(
		8080,
		storage.NewStore(redisClient),
		b.asyncEngine,
		b.codec,
		authenticator,
		b.modules,
		defaultAzureLocation,
		defaultAzureResourceGroup,
	)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Start starts all broker components (e.g. API server and async execution
// engine) and blocks until one of those components returns or fails.
func (b *broker) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errChan := make(chan error)
	// Start async engine
	go func() {
		select {
		case errChan <- &errAsyncEngineStopped{err: b.asyncEngine.Start(ctx)}:
		case <-ctx.Done():
		}
	}()
	// Start api server
	go func() {
		select {
		case errChan <- &errAPIServerStopped{err: b.apiServer.Start(ctx)}:
		case <-ctx.Done():
		}
	}()
	select {
	case <-ctx.Done():
		log.Debug("context canceled; broker shutting down")
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}
