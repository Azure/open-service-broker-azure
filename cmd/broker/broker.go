package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/api"
	apiFilters "github.com/Azure/open-service-broker-azure/pkg/api/filters"
	redisAsync "github.com/Azure/open-service-broker-azure/pkg/async/redis"
	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/broker"
	"github.com/Azure/open-service-broker-azure/pkg/crypto"
	"github.com/Azure/open-service-broker-azure/pkg/crypto/aes256"
	"github.com/Azure/open-service-broker-azure/pkg/crypto/noop"
	"github.com/Azure/open-service-broker-azure/pkg/http/filter"
	"github.com/Azure/open-service-broker-azure/pkg/http/filters"
	brokerLog "github.com/Azure/open-service-broker-azure/pkg/log"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/storage"
	"github.com/Azure/open-service-broker-azure/pkg/version"
	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
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
	logConfig, err := brokerLog.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(log.InfoLevel)
	logLevel := logConfig.GetLevel()
	log.WithField(
		"logLevel",
		strings.ToUpper(logLevel.String()),
	).Info("Setting log level")
	log.SetLevel(logLevel)

	azureConfig, err := azure.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize modules
	if err = initModules(azureConfig); err != nil {
		log.Fatal(err)
	}

	log.WithFields(
		log.Fields{
			"version": version.GetVersion(),
			"commit":  version.GetCommit(),
		},
	).Info("Open Service Broker for Azure starting")

	// Redis clients

	// Storage
	storageRedisConfig, err := storage.GetRedisConfig()
	if err != nil {
		log.Fatal(err)
	}
	storageRedisOpts := &redis.Options{
		Addr: fmt.Sprintf(
			"%s:%d",
			storageRedisConfig.GetHost(),
			storageRedisConfig.GetPort(),
		),
		Password:   storageRedisConfig.GetPassword(),
		DB:         storageRedisConfig.GetDB(),
		MaxRetries: 5,
	}
	if storageRedisConfig.IsTLSEnabled() {
		storageRedisOpts.TLSConfig = &tls.Config{
			ServerName: storageRedisConfig.GetHost(),
		}
	}
	storageRedisClient := redis.NewClient(storageRedisOpts)

	// Async
	asyncRedisConfig, err := redisAsync.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	asyncRedisOpts := &redis.Options{
		Addr: fmt.Sprintf(
			"%s:%d",
			asyncRedisConfig.GetHost(),
			asyncRedisConfig.GetPort(),
		),
		Password:   asyncRedisConfig.GetPassword(),
		DB:         asyncRedisConfig.GetDB(),
		MaxRetries: 5,
	}
	if asyncRedisConfig.IsTLSEnabled() {
		asyncRedisOpts.TLSConfig = &tls.Config{
			ServerName: asyncRedisConfig.GetHost(),
		}
	}
	asyncRedisClient := redis.NewClient(asyncRedisOpts)

	// Crypto
	cryptoConfig, err := crypto.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	var codec crypto.Codec
	cryptoScheme := cryptoConfig.GetEncryptionScheme()
	switch cryptoScheme {
	case crypto.AES256:
		codec, err = aes256.NewCodec([]byte(cryptoConfig.GetAES256Key()))
		if err != nil {
			log.Fatal(err)
		}
		log.WithField(
			"encryptionScheme",
			cryptoScheme,
		).Info("Sensitive instance and binding details will be encrypted")
	case crypto.NOOP:
		codec = noop.NewCodec()
		log.Warn(
			"ENCRYPTION IS DISABLED -- THIS IS NOT A SUITABLE OPTION FOR PRODUCTION",
		)
	}

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

	modulesConfig, err := service.GetModulesConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create broker
	broker, err := broker.NewBroker(
		storageRedisClient,
		asyncRedisClient,
		codec,
		filterChain,
		modules,
		modulesConfig.GetMinStability(),
		azureConfig.GetDefaultLocation(),
		azureConfig.GetDefaultResourceGroup(),
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
