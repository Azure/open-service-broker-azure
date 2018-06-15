package redis

import (
	"log"
	"os"
	"testing"

	"github.com/Azure/open-service-broker-azure/pkg/crypto"
	"github.com/Azure/open-service-broker-azure/pkg/crypto/noop"
)

func TestMain(m *testing.M) {
	if err := crypto.InitializeGlobalCodec(noop.NewCodec()); err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
