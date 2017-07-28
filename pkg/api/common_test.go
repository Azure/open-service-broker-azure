package api

import (
	"github.com/Azure/azure-service-broker/pkg/api/authenticator/always"
	fakeAsync "github.com/Azure/azure-service-broker/pkg/async/fake"
	"github.com/Azure/azure-service-broker/pkg/crypto/noop"
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/services/echo"
	memoryStorage "github.com/Azure/azure-service-broker/pkg/storage/memory"
	uuid "github.com/satori/go.uuid"
)

func getDisposableInstanceID() string {
	return uuid.NewV4().String()
}

func getDisposableServiceID() string {
	return uuid.NewV4().String()
}

func getDisposablePlanID() string {
	return uuid.NewV4().String()
}

func getDisposableBindingID() string {
	return uuid.NewV4().String()
}

func getTestServer() (*server, error) {
	echoModule := echo.New()
	echoCatalog, err := echoModule.GetCatalog()
	if err != nil {
		return nil, err
	}
	echoServices := echoCatalog.GetServices()
	echoServiceID := echoServices[0].GetID()
	modules := map[string]service.Module{
		echoServiceID: echoModule,
	}
	s, err := NewServer(
		8080,
		memoryStorage.NewStore(),
		fakeAsync.NewEngine(),
		noop.NewCodec(),
		always.NewAuthenticator(),
		modules,
	)
	if err != nil {
		return nil, err
	}
	return s.(*server), nil
}
