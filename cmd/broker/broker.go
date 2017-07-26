package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Azure/azure-service-broker/pkg/broker"
	"github.com/Azure/azure-service-broker/pkg/crypto/aes256"
	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

func main() {
	// Logging setup
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

	// Redis client
	redisConfig, err := getRedisConfig()
	if err != nil {
		log.Fatal(err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	// Crypto
	cryptoConfig, err := getCryptoConfig()
	if err != nil {
		log.Fatal(err)
	}
	codec, err := aes256.NewCodec(cryptoConfig.AES256Key)
	if err != nil {
		log.Fatal(err)
	}

	// Create broker
	broker, err := broker.NewBroker(redisClient, codec, modules)
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
