package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"open-service-broker-azure/pkg/api"
	apiFilters "open-service-broker-azure/pkg/api/filters"
	fakeAsync "open-service-broker-azure/pkg/async/fake"
	"open-service-broker-azure/pkg/crypto/noop"
	"open-service-broker-azure/pkg/http/filter"
	"open-service-broker-azure/pkg/http/filters"
	"open-service-broker-azure/pkg/services/fake"
	memoryStorage "open-service-broker-azure/pkg/storage/memory"

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

	filterChain := filter.NewChain(
		filters.NewBasicAuthFilter(username, password),
		apiFilters.NewAPIVersionFilter(),
	)

	noopCodec := noop.NewCodec()
	server, err := api.NewServer(
		8088,
		memoryStorage.NewStore(fakeCatalog, noopCodec),
		fakeAsync.NewEngine(),
		filterChain,
		fakeCatalog,
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
