package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/api"
	apiFilters "github.com/Azure/open-service-broker-azure/pkg/api/filters"
	async "github.com/Azure/open-service-broker-azure/pkg/async/redis"
	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/boot"
	"github.com/Azure/open-service-broker-azure/pkg/broker"
	"github.com/Azure/open-service-broker-azure/pkg/http/filter"
	"github.com/Azure/open-service-broker-azure/pkg/http/filters"
	brokerLog "github.com/Azure/open-service-broker-azure/pkg/log"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	storage "github.com/Azure/open-service-broker-azure/pkg/storage/redis"
	"github.com/Azure/open-service-broker-azure/pkg/version"

	log "github.com/Sirupsen/logrus"
)

func main() {
	// Initialize logging
	// Split log output across stdout and stderr, depending on severity
	// krancour: This functionality is currently dependent on a fork of
	// the github.com/Sirupsen/logrus package that lives in the split-streams
	// branch at github.com/krancour/logrus. (See Gopkg.toml)
	// We can resume using the upstream logrus if/when this PR is merged:
	// https://github.com/sirupsen/logrus/pull/671
	log.SetOutput(os.Stdout)
	log.SetErrOutput(os.Stderr)
	formatter := &log.TextFormatter{
		FullTimestamp: true,
	}
	log.SetFormatter(formatter)
	log.SetLevel(log.InfoLevel)

	log.WithFields(
		log.Fields{
			"version": version.GetVersion(),
			"commit":  version.GetCommit(),
		},
	).Info("Open Service Broker for Azure starting")

	logConfig, err := brokerLog.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	logLevel := logConfig.GetLevel()
	log.WithField(
		"logLevel",
		strings.ToUpper(logLevel.String()),
	).Info("Setting log level")
	log.SetLevel(logLevel)

	// Initialize catalog
	catalogConfig, err := service.GetCatalogConfigFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	azureConfig, err := azure.GetConfigFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	catalog, err := boot.GetCatalog(catalogConfig, azureConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Storage
	storageConfig, err := storage.GetConfigFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	store, err := storage.NewStore(catalog, storageConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Async
	asyncConfig, err := async.GetConfigFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	asyncEngine := async.NewEngine(asyncConfig)

	// Assemble the filter chain
	basicAuthConfig, err := api.GetBasicAuthConfig()
	if err != nil {
		log.Fatal(err)
	}
	filterChain := filter.NewChain(
		filters.NewBasicAuthFilter(
			basicAuthConfig.GetUsername(),
			basicAuthConfig.GetPassword(),
		),
		apiFilters.NewAPIVersionFilter(),
	)

	// Create broker
	broker, err := broker.NewBroker(
		store,
		asyncEngine,
		filterChain,
		catalog,
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

	// Run broker
	if err := broker.Start(ctx); err != nil {
		if err == ctx.Err() {
			// Allow some time for goroutines to shut down
			time.Sleep(time.Second * 3)
		} else {
			log.Fatal(err)
		}
	}
}
