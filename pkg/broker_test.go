package main

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"context"

	fakeAPI "github.com/Azure/azure-service-broker/pkg/api/fake"
	fakeAsync "github.com/Azure/azure-service-broker/pkg/async/fake"
	"github.com/stretchr/testify/assert"
)

func TestStartBrokerBlocks(t *testing.T) {
	s := fakeAPI.NewServer()
	e := fakeAsync.NewEngine()
	b := newBroker(s, e)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := b.start(ctx)
	assert.Equal(t, "context.deadlineExceededError", reflect.TypeOf(err).String())
}

func TestAPIServerErrorShutsDownBroker(t *testing.T) {
	s := fakeAPI.NewServer()
	someErr := errors.New("an error")
	s.RunBehavior = func() error {
		return someErr
	}
	e := fakeAsync.NewEngine()
	b := newBroker(s, e)
	err := b.start(context.Background())
	assert.Equal(t, someErr, err)
}

func TestAsyncEngineErrorShutsDownBroker(t *testing.T) {
	s := fakeAPI.NewServer()
	e := fakeAsync.NewEngine()
	someErr := errors.New("an error")
	e.RunBehavior = func() error {
		return someErr
	}
	b := newBroker(s, e)
	err := b.start(context.Background())
	assert.Equal(t, someErr, err)
}
