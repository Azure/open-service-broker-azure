package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/api"
	apiFilters "github.com/Azure/open-service-broker-azure/pkg/api/filters"
	"github.com/Azure/open-service-broker-azure/pkg/http/filter"
	"github.com/Azure/open-service-broker-azure/pkg/http/filters"
	"github.com/Azure/open-service-broker-azure/pkg/services/fake"
	memoryStorage "github.com/Azure/open-service-broker-azure/pkg/storage/memory"
	log "github.com/Sirupsen/logrus"
	fakeAsync "github.com/krancour/async/fake"
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

	server, err := api.NewServer(
		api.Config{
			Port: 8088,
		},
		memoryStorage.NewStore(fakeCatalog),
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

	if err := server.Run(ctx); err != nil {
		if err == ctx.Err() {
			// Allow some time for goroutines to shut down
			time.Sleep(time.Second * 3)
		} else {
			log.Fatal(err)
		}
	}
}
