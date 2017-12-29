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

	"github.com/Azure/open-service-broker-azure/pkg/api/authenticator/basic"
	"github.com/Azure/open-service-broker-azure/pkg/broker"
	"github.com/Azure/open-service-broker-azure/pkg/crypto/aes256"
	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

var (
	version string
	commit  string
)

func init() {
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
	logConfig, err := getLogConfig()
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(log.InfoLevel)
	log.WithField(
		"logLevel",
		strings.ToUpper(logConfig.Level.String()),
	).Info("setting log level")
	log.SetLevel(logConfig.Level)

	// Initialize modules
	if err = initModules(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.WithFields(
		log.Fields{
			"version": version,
			"commit":  commit,
		},
	).Info("Open Service Broker for Azure starting")

	// Redis client
	redisConfig, err := getRedisConfig()
	if err != nil {
		log.Fatal(err)
	}
	redisOpts := &redis.Options{
		Addr:       fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password:   redisConfig.Password,
		DB:         redisConfig.DB,
		MaxRetries: 5,
	}
	if redisConfig.EnableTLS {
		redisOpts.TLSConfig = &tls.Config{
			ServerName: redisConfig.Host,
		}
	}
	redisClient := redis.NewClient(redisOpts)

	// Crypto
	cryptoConfig, err := getCryptoConfig()
	if err != nil {
		log.Fatal(err)
	}
	codec, err := aes256.NewCodec([]byte(cryptoConfig.AES256Key))
	if err != nil {
		log.Fatal(err)
	}

	basicAuthConfig, err := getBasicAuthConfig()
	if err != nil {
		log.Fatal(err)
	}
	authenticator := basic.NewAuthenticator(
		basicAuthConfig.Username,
		basicAuthConfig.Password,
	)

	modulesConfig, err := getModulesConfig()
	if err != nil {
		log.Fatal(err)
	}

	azureConfig, err := getAzureConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create broker
	broker, err := broker.NewBroker(
		redisClient,
		codec,
		authenticator,
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
