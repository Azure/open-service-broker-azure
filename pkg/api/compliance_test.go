// +build compliance,!unit

package api

import (
	"context"
	"testing"

	"github.com/Azure/azure-service-broker/pkg/api/authenticator/basic"
	fakeAsync "github.com/Azure/azure-service-broker/pkg/async/fake"
	"github.com/Azure/azure-service-broker/pkg/crypto/noop"
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/services/fake"
	memoryStorage "github.com/Azure/azure-service-broker/pkg/storage/memory"
	log "github.com/Sirupsen/logrus"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
)

type basicAuthConfig struct {
	Username string `envconfig:"BASIC_AUTH_USERNAME" required:"true"`
	Password string `envconfig:"BASIC_AUTH_PASSWORD" required:"true"`
}

func getBasicAuthConfig() (basicAuthConfig, error) {
	bac := basicAuthConfig{}
	err := envconfig.Process("", &bac)
	return bac, err
}

func getComplianceTestServer() (*server, error) {
	fakeModule, err := fake.New()
	if err != nil {
		return nil, err
	}
	fakeCatalog, err := fakeModule.GetCatalog()
	if err != nil {
		return nil, err
	}
	modules := map[string]service.Module{
		fakeCatalog.GetServices()[0].GetID(): fakeModule,
	}

	basicAuthConfig, err := getBasicAuthConfig()
	if err != nil {
		log.Fatal(err)
	}
	authenticator := basic.NewAuthenticator(
		basicAuthConfig.Username,
		basicAuthConfig.Password,
	)

	s, err := NewServer(
		8080,
		memoryStorage.NewStore(),
		fakeAsync.NewEngine(),
		noop.NewCodec(),
		authenticator,
		modules,
	)
	if err != nil {
		return nil, err
	}
	return s.(*server), nil
}

//TestAPICompliance starts a test serverfor use with OSB api compliance testing
func TestAPICompliance(t *testing.T) {
	s, err := getComplianceTestServer()
	assert.Nil(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = s.listenAndServe(ctx)
	assert.Equal(t, ctx.Err(), err)
}
