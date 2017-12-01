package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Azure/azure-service-broker/pkg/api"
	"github.com/Azure/azure-service-broker/pkg/api/authenticator/basic"
	fakeAsync "github.com/Azure/azure-service-broker/pkg/async/fake"
	"github.com/Azure/azure-service-broker/pkg/crypto/noop"
	"github.com/Azure/azure-service-broker/pkg/services/fake"
	memoryStorage "github.com/Azure/azure-service-broker/pkg/storage/memory"
	log "github.com/Sirupsen/logrus"
)

func main() {
	fakeModule, err := fake.New()
	if err != nil {
		log.Fatal(err)
	}
	fakeCatalog, err := fakeModule.GetCatalog()
	if err != nil {
		log.Fatal(err)
	}

	username := "username"
	password := "password"

	authenticator := basic.NewAuthenticator(
		username,
		password,
	)

	server, err := api.NewServer(
		8080,
		memoryStorage.NewStore(),
		fakeAsync.NewEngine(),
		noop.NewCodec(),
		authenticator,
		fakeCatalog,
		" ",
		" ",
	)

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		signal := <-sigChan
		log.WithField(
			"signal",
			signal,
		).Debug("signal received; shutting down")
		cancel()
	}()

	if err := server.Start(ctx); err != nil {
		if err == ctx.Err() {
			// Allow some time for goroutines to shut down
			time.Sleep(time.Second * 3)
		} else {
			log.Fatal(err)
		}
	}
}
