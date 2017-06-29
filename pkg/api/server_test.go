package api

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"time"

	"errors"

	fakeAsync "github.com/Azure/azure-service-broker/pkg/async/fake"
	"github.com/Azure/azure-service-broker/pkg/service"
	fakeStorage "github.com/Azure/azure-service-broker/pkg/storage/fake"
	"github.com/stretchr/testify/assert"
)

var (
	store       = fakeStorage.NewStore()
	asyncEngine = fakeAsync.NewEngine()
)

func TestStartServerBlocks(t *testing.T) {
	s, err := NewServer(8080, store, asyncEngine, []service.Module{})
	assert.Nil(t, err)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = s.Start(ctx)
	assert.Equal(t, "context.deadlineExceededError", reflect.TypeOf(err).String())
}

func TestErrorShutsDownAPIServer(t *testing.T) {
	s, err := NewServer(8080, store, asyncEngine, []service.Module{})
	assert.Nil(t, err)
	server := s.(*server)
	someErr := errors.New("an error")
	server.listenAndServe = func(string, http.Handler) error {
		return someErr
	}
	err = s.Start(context.Background())
	assert.Equal(t, someErr, err)
}
