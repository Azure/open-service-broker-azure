package broker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/http/filter"
	"github.com/Azure/open-service-broker-azure/pkg/service"

	fakeAPI "github.com/Azure/open-service-broker-azure/pkg/api/fake"
	fakeAsync "github.com/Azure/open-service-broker-azure/pkg/async/fake"
	"github.com/stretchr/testify/assert"
)

var errSome = errors.New("an error")

func TestBrokerStartBlocksUntilAsyncEngineErrors(t *testing.T) {
	apiServerStopped := false
	svr := fakeAPI.NewServer()
	svr.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		apiServerStopped = true
		return ctx.Err()
	}
	e := fakeAsync.NewEngine()
	e.RunBehavior = func(context.Context) error {
		return errSome
	}
	b, err := getTestBroker()
	assert.Nil(t, err)
	b.asyncEngine = e
	b.apiServer = svr
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = b.Start(ctx)
	assert.Equal(t, &errAsyncEngineStopped{err: errSome}, err)
	time.Sleep(time.Second)
	assert.True(t, apiServerStopped)
}

func TestBrokerStartBlocksUntilAsyncEngineReturns(t *testing.T) {
	apiServerStopped := false
	svr := fakeAPI.NewServer()
	svr.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		apiServerStopped = true
		return ctx.Err()
	}
	e := fakeAsync.NewEngine()
	e.RunBehavior = func(context.Context) error {
		return nil
	}
	b, err := getTestBroker()
	assert.Nil(t, err)
	b.asyncEngine = e
	b.apiServer = svr
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = b.Start(ctx)
	assert.Equal(t, &errAsyncEngineStopped{}, err)
	time.Sleep(time.Second)
	assert.True(t, apiServerStopped)
}

func TestBrokerStartBlocksUntilAPIServerErrors(t *testing.T) {
	svr := fakeAPI.NewServer()
	svr.RunBehavior = func(context.Context) error {
		return errSome
	}
	asyncEngineStopped := false
	e := fakeAsync.NewEngine()
	e.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		asyncEngineStopped = true
		return ctx.Err()
	}
	b, err := getTestBroker()
	assert.Nil(t, err)
	b.asyncEngine = e
	b.apiServer = svr
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = b.Start(ctx)
	assert.Equal(t, &errAPIServerStopped{err: errSome}, err)
	time.Sleep(time.Second)
	assert.True(t, asyncEngineStopped)
}

func TestBrokerStartBlocksUntilAPIServerReturns(t *testing.T) {
	svr := fakeAPI.NewServer()
	svr.RunBehavior = func(context.Context) error {
		return nil
	}
	asyncEngineStopped := false
	e := fakeAsync.NewEngine()
	e.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		asyncEngineStopped = true
		return ctx.Err()
	}
	b, err := getTestBroker()
	assert.Nil(t, err)
	b.asyncEngine = e
	b.apiServer = svr
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = b.Start(ctx)
	assert.Equal(t, &errAPIServerStopped{}, err)
	time.Sleep(time.Second)
	assert.True(t, asyncEngineStopped)
}

func TestBrokerStartBlocksUntilContextCanceled(t *testing.T) {
	apiServerStopped := false
	svr := fakeAPI.NewServer()
	svr.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		apiServerStopped = true
		return ctx.Err()
	}
	asyncEngineStopped := false
	e := fakeAsync.NewEngine()
	e.RunBehavior = func(ctx context.Context) error {
		<-ctx.Done()
		asyncEngineStopped = true
		return ctx.Err()
	}
	b, err := getTestBroker()
	assert.Nil(t, err)
	b.asyncEngine = e
	b.apiServer = svr
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = b.Start(ctx)
	assert.Equal(t, ctx.Err(), err)
	time.Sleep(time.Second)
	assert.True(t, apiServerStopped)
	assert.True(t, asyncEngineStopped)
}

func getTestBroker() (*broker, error) {
	b, err := NewBroker(
		nil,
		fakeAsync.NewEngine(),
		filter.NewChain(),
		service.NewCatalog(nil),
		azure.NewConfigWithDefaults(),
	)
	if err != nil {
		return nil, err
	}
	return b.(*broker), nil
}
