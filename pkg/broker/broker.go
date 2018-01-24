package broker

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/api"
	"github.com/Azure/open-service-broker-azure/pkg/api/filter"
	"github.com/Azure/open-service-broker-azure/pkg/async"
	redisAsync "github.com/Azure/open-service-broker-azure/pkg/async/redis"
	"github.com/Azure/open-service-broker-azure/pkg/crypto"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/storage"
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
	catalog     service.Catalog
}

// NewBroker returns a new Broker
func NewBroker(
	storageRedisClient *redis.Client,
	asyncRedisClient *redis.Client,
	codec crypto.Codec,
	filterChain filter.Filter,
	modules []service.Module,
	minStability service.Stability,
	defaultAzureLocation string,
	defaultAzureResourceGroup string,
) (Broker, error) {
	// Consolidate the catalogs from all the individual modules into a single
	// catalog. Check as we go along to make sure that no two modules provide
	// services having the same ID.
	services := []service.Service{}
	usedServiceIDs := map[string]string{}
	for _, module := range modules {
		if module.GetStability() >= minStability {
			moduleName := module.GetName()
			catalog, err := module.GetCatalog()
			if err != nil {
				return nil, fmt.Errorf(
					`error retrieving catalog from module "%s": %s`,
					moduleName,
					err,
				)
			}
			for _, svc := range catalog.GetServices() {
				serviceID := svc.GetID()
				if moduleNameForUsedServiceID, ok := usedServiceIDs[serviceID]; ok {
					return nil, fmt.Errorf(
						`modules "%s" and "%s" both provide a service with the id "%s"`,
						moduleNameForUsedServiceID,
						moduleName,
						serviceID,
					)
				}
				services = append(services, svc)
				usedServiceIDs[serviceID] = moduleName
			}
		}
	}
	catalog := service.NewCatalog(services)
	b := &broker{
		store:       storage.NewStore(storageRedisClient, catalog, codec),
		asyncEngine: redisAsync.NewEngine(asyncRedisClient),
		catalog:     catalog,
	}

	err := b.asyncEngine.RegisterJob(
		"executeProvisioningStep",
		b.executeProvisioningStep,
	)
	if err != nil {
		return nil, errors.New(
			"error registering async job for executing provisioning steps",
		)
	}
	err = b.asyncEngine.RegisterJob("executeUpdatingStep", b.executeUpdatingStep)
	if err != nil {
		return nil, errors.New(
			"error registering async job for executing updating steps",
		)
	}
	err = b.asyncEngine.RegisterJob(
		"executeDeprovisioningStep",
		b.executeDeprovisioningStep,
	)
	if err != nil {
		return nil, errors.New(
			"error registering async job for executing deprovisioning steps",
		)
	}
	b.apiServer, err = api.NewServer(
		8080,
		b.store,
		b.asyncEngine,
		filterChain,
		b.catalog,
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
		case errChan <- &errAsyncEngineStopped{err: b.asyncEngine.Run(ctx)}:
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
