package main

import (
	"context"

	"github.com/Azure/azure-service-broker/pkg/api"
	"github.com/Azure/azure-service-broker/pkg/async"
)

type broker struct {
	apiServer   api.Server
	asyncEngine async.Engine
}

func newBroker(
	apiServer api.Server,
	asyncEngine async.Engine,
) *broker {
	return &broker{
		apiServer:   apiServer,
		asyncEngine: asyncEngine,
	}
}

func (b *broker) start(ctx context.Context) error {
	errChan := make(chan error)
	defer close(errChan)
	// Start async machinery
	go func() {
		errChan <- b.asyncEngine.Start(ctx)
	}()
	// Start api server
	go func() {
		errChan <- b.apiServer.Start(ctx)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}
