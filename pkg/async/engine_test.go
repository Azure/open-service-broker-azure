package async

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	machinery "github.com/Azure/azure-service-broker/pkg/machinery"
	fakeMachinery "github.com/Azure/azure-service-broker/pkg/machinery/fake"
	"github.com/Azure/azure-service-broker/pkg/service"
	fakeStorage "github.com/Azure/azure-service-broker/pkg/storage/fake"
	"github.com/stretchr/testify/assert"
)

var (
	s = fakeStorage.NewStore()
	m = fakeMachinery.NewServer()
)

func TestStartEngineBlocks(t *testing.T) {
	e, err := NewEngine(s, m, []service.Module{})
	assert.Nil(t, err)
	eng := e.(*engine)
	eng.getWorker = func(machinery.Server) machinery.Worker {
		return fakeMachinery.NewWorker()
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err = e.Start(ctx)
	assert.Equal(t, "context.deadlineExceededError", reflect.TypeOf(err).String())
}

func TestErrorShutsDownEngine(t *testing.T) {
	e, err := NewEngine(s, m, []service.Module{})
	assert.Nil(t, err)
	eng := e.(*engine)
	someErr := errors.New("an error")
	eng.getWorker = func(machinery.Server) machinery.Worker {
		worker := fakeMachinery.NewWorker()
		worker.RunBehavior = func() error {
			return someErr
		}
		return worker
	}
	err = e.Start(context.Background())
	assert.Equal(t, someErr, err)
}
