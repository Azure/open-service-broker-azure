package main

import (
	"fmt"
	"log"

	"context"

	"github.com/Azure/azure-service-broker/pkg/api"
	"github.com/Azure/azure-service-broker/pkg/async"
	"github.com/Azure/azure-service-broker/pkg/storage"
	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/go-redis/redis"
)

func main() {
	// Get Redis config-- we'll use this for both the storage layer and for
	// the async machinery
	redisConfig, err := getRedisConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Get a Redis-based implementation of the storage.Store interface
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password, // "" == no password
		DB:       redisConfig.DB,
	})
	store := storage.NewStore(redisClient)

	// Get the underlying machinery for the async engine
	redisURL := fmt.Sprintf(
		"redis://%s@%s:%d/%d",
		redisConfig.Password,
		redisConfig.Host,
		redisConfig.Port,
		redisConfig.DB,
	)
	var machineryConfig = &config.Config{
		Broker:        redisURL,
		ResultBackend: redisURL,
	}
	machineryServer, err := machinery.NewServer(machineryConfig)
	if err != nil {
		log.Fatal(err)
	}
	// And get the async engine
	asyncEngine, err := async.NewEngine(store, machineryServer, getModules())
	if err != nil {
		log.Fatal(err)
	}

	// Get web config and an apiServer
	webConfig, err := getWebConfig()
	if err != nil {
		log.Fatal(err)
	}
	apiServer, err := api.NewServer(
		webConfig.Port,
		store,
		asyncEngine,
		getModules(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Set up and start the broker
	broker := newBroker(apiServer, asyncEngine)
	// Start the broker
	if err := broker.start(context.Background()); err != nil {
		log.Fatal(err)
	}
}
