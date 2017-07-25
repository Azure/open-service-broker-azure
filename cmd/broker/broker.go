package main

import (
	"fmt"

	"context"

	"github.com/Azure/azure-service-broker/pkg/broker"
	"github.com/Azure/azure-service-broker/pkg/crypto/aes256"
	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

func main() {
	redisConfig, err := getRedisConfig()
	if err != nil {
		log.Fatal(err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})
	cryptoConfig, err := getCryptoConfig()
	if err != nil {
		log.Fatal(err)
	}
	codec, err := aes256.NewCodec(cryptoConfig.AES256Key)
	if err != nil {
		log.Fatal(err)
	}
	broker, err := broker.NewBroker(redisClient, codec, modules)
	if err != nil {
		log.Fatal(err)
	}
	if err := broker.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}
