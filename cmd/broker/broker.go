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

	apiFilters "github.com/Azure/open-service-broker-azure/pkg/api/filters"
	"github.com/Azure/open-service-broker-azure/pkg/broker"
	"github.com/Azure/open-service-broker-azure/pkg/config"
	"github.com/Azure/open-service-broker-azure/pkg/crypto/aes256"
	"github.com/Azure/open-service-broker-azure/pkg/http/filter"
	"github.com/Azure/open-service-broker-azure/pkg/http/filters"
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
	logConfig, err := config.GetLogConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(log.InfoLevel)
	logLevel := logConfig.GetLevel()
	log.WithField(
		"logLevel",
		strings.ToUpper(logLevel.String()),
	).Info("setting log level")
	log.SetLevel(logLevel)

	azureConfig, err := config.GetAzureConfig()
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
	redisConfig, err := config.GetRedisConfig()
	if err != nil {
		log.Fatal(err)
	}
	storageRedisOpts := &redis.Options{
		Addr: fmt.Sprintf(
			"%s:%d",
			redisConfig.GetHost(),
			redisConfig.GetPort(),
		),
		Password:   redisConfig.GetPassword(),
		DB:         redisConfig.GetStorageDB(),
		MaxRetries: 5,
	}
	asyncRedisOpts := &redis.Options{
		Addr: fmt.Sprintf(
			"%s:%d",
			redisConfig.GetHost(),
			redisConfig.GetPort(),
		),
		Password:   redisConfig.GetPassword(),
		DB:         redisConfig.GetAsyncDB(),
		MaxRetries: 5,
	}
	if redisConfig.IsTLSEnabled() {
		storageRedisOpts.TLSConfig = &tls.Config{
			ServerName: redisConfig.GetHost(),
		}
		asyncRedisOpts.TLSConfig = &tls.Config{
			ServerName: redisConfig.GetHost(),
		}
	}
	storageRedisClient := redis.NewClient(storageRedisOpts)
	asyncRedisClient := redis.NewClient(asyncRedisOpts)

	// Crypto
	cryptoConfig, err := config.GetCryptoConfig()
	if err != nil {
		log.Fatal(err)
	}
	codec, err := aes256.NewCodec([]byte(cryptoConfig.GetAES256Key()))
	if err != nil {
		log.Fatal(err)
	}

	// Assemble the filter chain
	basicAuthConfig, err := config.GetBasicAuthConfig()
	if err != nil {
		log.Fatal(err)
	}
	filterChain := filter.NewChain(
		filters.NewBasicAuthFilter(
			basicAuthConfig.Username,
			basicAuthConfig.Password,
		),
		apiFilters.NewAPIVersionFilter(),
	)

	modulesConfig, err := config.GetModulesConfig()
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
		modulesConfig.MinStability,
		azureConfig.DefaultLocation,
		azureConfig.DefaultResourceGroup,
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
