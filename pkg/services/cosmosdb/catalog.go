package cosmosdb

import "github.com/Azure/open-service-broker-azure/pkg/service"

const kindKey = "kind"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
			service.NewService(
				&service.ServiceProperties{
					ID:   "6330de6f-a561-43ea-a15e-b99f44d183e6",
					Name: "azure-cosmos-db",
					Description: "Azure Cosmos DB is a globally distributed, " +
						"multi-model database service accessible via SQL, " +
						"Gremlin (Graph), and Table (Key-Value) APIs",
					Metadata: &service.ServiceMetadata{
						DisplayName: "Azure Cosmos DB",
						ImageUrl: "https://azure.microsoft.com/svghandler/cosmos-db/" +
							"?width=200",
						LongDescription: "Globally distributed, multi-model database service" +
							" (Experimental)",
						DocumentationUrl: "https://docs.microsoft.com/en-us/azure/cosmos-db/",
						SupportUrl:       "https://azure.microsoft.com/en-us/support/",
					},
					Bindable: true,
					Tags: []string{"Azure",
						"CosmosDB",
						"Database",
						"SQL",
						"Gremlin",
						"Graph",
						"Table",
						"Key-Value",
					},
				},
				m.cosmosAccountManager,
				service.NewPlan(&service.PlanProperties{
					ID:          "71168d1a-c704-49ff-8c79-214dd3d6f8eb",
					Name:        "cosmos-db-sql-api",
					Description: "Creates a Database Account with the SQL API",
					Free:        false,
					Extended: map[string]interface{}{
						kindKey: databaseKindGlobalDocumentDB,
					},
					Metadata: &service.ServicePlanMetadata{
						DisplayName: "Azure DocumentDB",
					},
				}),
			),
			service.NewService(
				&service.ServiceProperties{
					ID:          "8797a079-5346-4e84-8018-b7d5ea5c0e3a",
					Name:        "azure-cosmos-mongo-db",
					Description: "MongoDB on Azure (Experimental) provided by CosmosDB",
					Metadata: &service.ServiceMetadata{
						DisplayName: "Azure Cosmos DB (MongoDB)",
						ImageUrl: "https://azure.microsoft.com/svghandler/cosmos-db/" +
							"?width=200",
						LongDescription: "Globally distributed, multi-model database service" +
							" (Experimental)",
						DocumentationUrl: "https://docs.microsoft.com/en-us/azure/cosmos-db/",
						SupportUrl:       "https://azure.microsoft.com/en-us/support/",
					},
					Bindable: true,
					Tags: []string{"Azure",
						"CosmosDB",
						"Database",
						"MongoDB",
					},
				},
				m.mongoAccountManager,
				service.NewPlan(&service.PlanProperties{
					ID:          "86fdda05-78d7-4026-a443-1325928e7b02",
					Name:        "mongo-db",
					Description: "Creates a Database Account with the MongoDB API",
					Free:        false,
					Extended: map[string]interface{}{
						kindKey: databaseKindMongoDB,
					},
					Metadata: &service.ServicePlanMetadata{
						DisplayName: "Azure MongoDB",
					},
				}),
			),
		}),
		nil
}
