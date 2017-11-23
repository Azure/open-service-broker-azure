package api

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	errSome = errors.New("an error")
)

func TestServerStartBlocksUntilListenAndServeErrors(t *testing.T) {
	s, _, err := getTestServer("", "")
	assert.Nil(t, err)
	s.listenAndServe = func(context.Context) error {
		return errSome
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = s.Start(ctx)
	assert.Equal(t, &errHTTPServerStopped{err: errSome}, err)
}

func TestServerStartBlocksUntilListenAndServeReturns(t *testing.T) {
	s, _, err := getTestServer("", "")
	assert.Nil(t, err)
	s.listenAndServe = func(context.Context) error {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = s.Start(ctx)
	assert.Equal(t, &errHTTPServerStopped{}, err)
}

func TestServerStartBlocksUntilContextCanceled(t *testing.T) {
	s, _, err := getTestServer("", "")
	assert.Nil(t, err)
	s.listenAndServe = func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = s.Start(ctx)
	assert.Equal(t, ctx.Err(), err)
}

func TestServerListenAndServeBlocksUntilContextCanceled(t *testing.T) {
	s, _, err := getTestServer("", "")
	assert.Nil(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = s.listenAndServe(ctx)
	assert.Equal(t, ctx.Err(), err)
}
