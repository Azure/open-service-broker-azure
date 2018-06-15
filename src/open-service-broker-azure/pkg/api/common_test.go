package api

import (
	"fmt"

	fakeAsync "open-service-broker-azure/pkg/async/fake"
	"open-service-broker-azure/pkg/crypto/noop"
	"open-service-broker-azure/pkg/http/filter"
	"open-service-broker-azure/pkg/services/fake"
	memoryStorage "open-service-broker-azure/pkg/storage/memory"
	uuid "github.com/satori/go.uuid"
)

type ArbitraryType struct {
	Foo string `json:"foo"`
}

const fooValue = "bar"

var (
	testArbitraryObject = &ArbitraryType{
		Foo: fooValue,
	}
	testArbitraryObjectJSON = []byte(fmt.Sprintf(`{"foo":"%s"}`, fooValue))
	testArbitraryMap        = map[string]interface{}{
		"foo": "bar",
	}
	testArbitraryMapJSON = []byte(fmt.Sprintf(`{"foo":"%s"}`, fooValue))
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
	s, err := NewServer(
		8080,
		memoryStorage.NewStore(fakeCatalog, noop.NewCodec()),
		fakeAsync.NewEngine(),
		filter.NewChain(),
		fakeCatalog,
	)
	if err != nil {
		return nil, nil, err
	}
	return s.(*server), fakeModule, nil
}
