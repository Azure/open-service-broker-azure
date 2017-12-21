package main

import (
	"github.com/Azure/open-service-broker-azure/pkg/broker"
	"github.com/Azure/open-service-broker-azure/pkg/crypto"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/storage"
	"github.com/Azure/open-service-broker-azure/pkg/storage/cosmosdb"
	"github.com/go-redis/redis"
	mgo "gopkg.in/mgo.v2"
)

// getCreateStorageFunc creates a new broker.CreateStorageFunc that creates
// and returns the proper storage driver based on the given storageConfig.
//
// The only supported storage drivers are 'redis' and 'cosmosdb'. This function
// will return a StorageFunc for redis if storageConfig.StorageType is not
// 'cosmosdb'
func getCreateStorageFunc(
	redisClient *redis.Client,
	storageConfig storageConfig,
	cosmosConfig cosmosDBConfig,
) broker.CreateStorageFunc {
	return func(catalog service.Catalog, codec crypto.Codec) (storage.Store, error) {
		if storageConfig.StorageType == StorageTypeCosmosDB {
			db, err := getMongoDB(cosmosConfig.ConnectionURL, cosmosConfig.DBName)
			if err != nil {
				return nil, err
			}
			return cosmosdb.NewStore(
				db,
				cosmosConfig.InstanceCollectionName,
				cosmosConfig.BindingCollectionName,
			), nil
		}
		return storage.NewStore(redisClient, catalog, codec), nil
	}
}

func getMongoDB(url, dbName string) (*mgo.Database, error) {
	sess, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	return sess.DB(dbName), nil
}
