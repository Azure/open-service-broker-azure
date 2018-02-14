// +build !unit

package lifecycle

import (
	"crypto/tls"
	"fmt"
	"net"
	"strconv"
	"time"

	cosmosSDK "github.com/Azure/azure-sdk-for-go/arm/cosmos-db"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
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
			module:                 cosmosdb.New(armDeployer, dbAccountsClient),
			description:            "DocumentDB",
			serviceID:              "6330de6f-a561-43ea-a15e-b99f44d183e6",
			planID:                 "71168d1a-c704-49ff-8c79-214dd3d6f8eb",
			location:               "eastus",
			provisioningParameters: &cosmosdb.ProvisioningParameters{},
			bindingParameters:      &cosmosdb.BindingParameters{},
			testCredentials:        testDocumentDBCreds(),
		},
		{ // MongoDB
			module:                 cosmosdb.New(armDeployer, dbAccountsClient),
			description:            "MongoDB",
			serviceID:              "8797a079-5346-4e84-8018-b7d5ea5c0e3a",
			planID:                 "86fdda05-78d7-4026-a443-1325928e7b02",
			location:               "southcentralus",
			provisioningParameters: &cosmosdb.ProvisioningParameters{},
			bindingParameters:      &cosmosdb.BindingParameters{},
			testCredentials:        testMongoDBCreds(),
		},
	}, nil
}

func testDocumentDBCreds() func(credentials service.Credentials) error {
	return func(credentials service.Credentials) error {
		//cdts, err := convertToCosmosdbCredentials(credentials)
		//if err != nil {
		//	return err
		//}

		// Found no usable go package.
		// Tried github.com/nerdylikeme/go-documentdb.
		// It is out-of-date. Got following error:
		// {Code:"BadRequest", Message:"Invalid API version.
		//   Ensure a valid x-ms-version header value is passed.\r\nActivityId:
		//   7073673c-498a-4094-8749-5f5b3d5b838a"}
		// TODO:
		//   opt1. use REST API to request:
		//     https://docs.microsoft.com/en-us/rest/api/documentdb
		//   opt2. track this repo:
		//     https://github.com/a8m/documentdb-go
		return nil
	}
}

func testMongoDBCreds() func(credentials service.Credentials) error {
	return func(credentials service.Credentials) error {
		cdts, ok := credentials.(*cosmosdb.Credentials)
		if !ok {
			return fmt.Errorf("error casting credentials as *cosmosdb.Credentials")
		}

		// The following process bases on
		// https://docs.microsoft.com/en-us/azure/cosmos-db/create-mongodb-golang

		// DialInfo holds options for establishing a session with a MongoDB cluster.
		dialInfo := &mgo.DialInfo{
			Addrs:    []string{cdts.Host + ":" + strconv.Itoa(cdts.Port)},
			Timeout:  60 * time.Second,
			Database: "database",
			Username: cdts.Username,
			Password: cdts.Password,
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
