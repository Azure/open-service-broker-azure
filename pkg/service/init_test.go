package service

import (
	"log"
	"os"
	"testing"

	"github.com/Azure/open-service-broker-azure/pkg/crypto"
	"github.com/Azure/open-service-broker-azure/pkg/crypto/aes256"
)

func TestMain(m *testing.M) {
	codec, err := aes256.NewCodec(
		aes256.Config{
			Key: "AES256Key-32Characters1234567890",
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := crypto.InitializeGlobalCodec(codec); err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}
