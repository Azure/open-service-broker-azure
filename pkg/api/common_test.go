package api

import (
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/api/authenticator/always"
	fakeAsync "github.com/Azure/azure-service-broker/pkg/async/fake"
	"github.com/Azure/azure-service-broker/pkg/crypto/noop"
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/services/fake"
	memoryStorage "github.com/Azure/azure-service-broker/pkg/storage/memory"
	uuid "github.com/satori/go.uuid"
)

type ArbitraryType struct {
	Foo string `json:"foo"`
}

// SetResourceGroup is implemented only so that ArbitraryType will fulfill
// the ProvisioningParameters interface. This function isn't used anywhere.
func (a *ArbitraryType) SetResourceGroup(resourceGroup string) {}

const fooValue = "bar"

var (
	testArbitraryObject = &ArbitraryType{
		Foo: fooValue,
	}
	testArbitraryObjectJSON = []byte(fmt.Sprintf(`{"foo":"%s"}`, fooValue))
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

func getTestServer() (*server, *fake.Module, error) {
	fakeModule, err := fake.New()
	if err != nil {
		return nil, nil, err
	}
	fakeCatalog, err := fakeModule.GetCatalog()
	if err != nil {
		return nil, nil, err
	}
	modules := map[string]service.Module{
		fakeCatalog.GetServices()[0].GetID(): fakeModule,
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
		return nil, nil, err
	}
	return s.(*server), fakeModule, nil
}
