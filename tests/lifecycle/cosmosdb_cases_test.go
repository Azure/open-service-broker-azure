// +build !unit

package lifecycle

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	cosmosSDK "github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb" // nolint: lll
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/services/cosmosdb"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func getCosmosdbCases(
	azureEnvironment azure.Environment,
	subscriptionID string,
	authorizer autorest.Authorizer,
	armDeployer arm.Deployer,
) ([]serviceLifecycleTestCase, error) {
	dbAccountsClient := cosmosSDK.NewDatabaseAccountsClientWithBaseURI(
		azureEnvironment.ResourceManagerEndpoint,
		subscriptionID,
	)
	dbAccountsClient.Authorizer = authorizer
	return []serviceLifecycleTestCase{
		{ // DocumentDB
			module:      cosmosdb.New(armDeployer, dbAccountsClient),
			description: "CosmosDB",
			serviceID:   "6330de6f-a561-43ea-a15e-b99f44d183e6",
			planID:      "71168d1a-c704-49ff-8c79-214dd3d6f8eb",
			location:    "eastus",
		},
		{ // MongoDB
			module:          cosmosdb.New(armDeployer, dbAccountsClient),
			description:     "MongoDB API on CosmosDB",
			serviceID:       "8797a079-5346-4e84-8018-b7d5ea5c0e3a",
			planID:          "86fdda05-78d7-4026-a443-1325928e7b02",
			location:        "southcentralus",
			testCredentials: testMongoDBCreds(),
		},
	}, nil
}

func testMongoDBCreds() func(credentials map[string]interface{}) error {
	return func(credentials map[string]interface{}) error {
		// The following process is based on
		// https://docs.microsoft.com/en-us/azure/cosmos-db/create-mongodb-golang

		// DialInfo holds options for establishing a session with a MongoDB cluster.
		dialInfo := &mgo.DialInfo{
			Addrs: []string{
				fmt.Sprintf(
					"%s:%d",
					credentials["host"].(string),
					int(credentials["port"].(float64)),
				),
			},
			Timeout:  60 * time.Second,
			Database: "database",
			Username: credentials["username"].(string),
			Password: credentials["password"].(string),
			DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
				return tls.Dial("tcp", addr.String(), &tls.Config{})
			},
		}

		// Create a session which maintains a pool of socket connections
		// to our Azure Cosmos DB MongoDB database.
		session, err := mgo.DialWithInfo(dialInfo)

		if err != nil {
			return err
		}

		defer session.Close()

		session.SetSafe(&mgo.Safe{})

		collection := session.DB("database").C("package")

		// Model
		type Package struct {
			ID            bson.ObjectId `bson:"_id,omitempty"`
			FullName      string
			Description   string
			StarsCount    int
			ForksCount    int
			LastUpdatedBy string
		}

		// insert Document in collection
		err = collection.Insert(&Package{
			FullName:      "react",
			Description:   "A framework for building native apps with React.",
			ForksCount:    11392,
			StarsCount:    48794,
			LastUpdatedBy: "shergin",
		})

		return err
	}
}
