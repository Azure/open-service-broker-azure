package api

import (
	"context"
	"testing"
	"time"

	"errors"

	fakeAsync "github.com/Azure/azure-service-broker/pkg/async/fake"
	fakeStorage "github.com/Azure/azure-service-broker/pkg/storage/fake"
	"github.com/stretchr/testify/assert"
)

var (
	store       = fakeStorage.NewStore()
	asyncEngine = fakeAsync.NewEngine()
	errSome     = errors.New("an error")
)

func TestServerStartBlocksUntilListenAndServeErrors(t *testing.T) {
	s, err := NewServer(8080, store, asyncEngine, nil, nil, nil, nil)
	assert.Nil(t, err)
	server := s.(*server)
	server.listenAndServe = func(context.Context) error {
		return errSome
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = s.Start(ctx)
	assert.Equal(t, &errHTTPServerStopped{err: errSome}, err)
}

func TestServerStartBlocksUntilListenAndServeReturns(t *testing.T) {
	s, err := NewServer(8080, store, asyncEngine, nil, nil, nil, nil)
	assert.Nil(t, err)
	server := s.(*server)
	server.listenAndServe = func(context.Context) error {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = s.Start(ctx)
	assert.Equal(t, &errHTTPServerStopped{}, err)
}

func TestServerStartBlocksUntilContextCanceled(t *testing.T) {
	s, err := NewServer(8080, store, asyncEngine, nil, nil, nil, nil)
	assert.Nil(t, err)
	server := s.(*server)
	server.listenAndServe = func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = s.Start(ctx)
	assert.Equal(t, ctx.Err(), err)
}

func TestServerListenAndServeBlocksUntilContextCanceled(t *testing.T) {
	s, err := NewServer(8080, store, asyncEngine, nil, nil, nil, nil)
	server := s.(*server)
	assert.Nil(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = server.listenAndServe(ctx)
	assert.Equal(t, ctx.Err(), err)
}
