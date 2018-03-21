package cosmosdb

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
	commonSchema := &service.ParameterSchemas{
		ServiceInstances: &service.InstanceSchema{
			Create: &service.InputParameters{
				Parameters: service.GetCommonSchema(),
			},
		},
	}
	return service.NewCatalog([]service.Service{
			service.NewService(
				&service.ServiceProperties{
					ID:          "6330de6f-a561-43ea-a15e-b99f44d183e6",
					Name:        "azure-cosmosdb-sql-account",
					Description: "Azure Cosmos DB Database Account (SQL API)",
					Metadata: &service.ServiceMetadata{
						DisplayName: "Azure Cosmos DB",
						ImageURL: "https://azure.microsoft.com/svghandler/cosmos-db/" +
							"?width=200",
						LongDescription: "Globally distributed, multi-model database service" +
							" (Experimental).",
						DocumentationURL: "https://docs.microsoft.com/en-us/azure/cosmos-db/",
						SupportURL:       "https://azure.microsoft.com/en-us/support/",
					},
					Bindable: true,
					Tags: []string{"Azure",
						"CosmosDB",
						"Database",
						"SQL",
					},
				},
				m.cosmosAccountManager,
				service.NewPlan(&service.PlanProperties{
					ID:          "71168d1a-c704-49ff-8c79-214dd3d6f8eb",
					Name:        "account",
					Description: "Database Account with the SQL API",
					Free:        false,
					Metadata: &service.ServicePlanMetadata{
						DisplayName: "Azure CosmosDB (SQL API)",
					},
					ParameterSchemas: commonSchema,
				}),
			),
			service.NewService(
				&service.ServiceProperties{
					ID:          "8797a079-5346-4e84-8018-b7d5ea5c0e3a",
					Name:        "azure-cosmosdb-mongo-account",
					Description: "Azure Cosmos DB Database Account (MongoDB API)",
					Metadata: &service.ServiceMetadata{
						DisplayName: "Azure Cosmos DB (MongoDB)",
						ImageURL: "https://azure.microsoft.com/svghandler/cosmos-db/" +
							"?width=200",
						LongDescription: "Globally distributed, multi-model database service" +
							" (Experimental)",
						DocumentationURL: "https://docs.microsoft.com/en-us/azure/cosmos-db/",
						SupportURL:       "https://azure.microsoft.com/en-us/support/",
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
					Name:        "account",
					Description: "Database Account with the MongoDB API",
					Free:        false,
					Metadata: &service.ServicePlanMetadata{
						DisplayName: "Azure MongoDB",
					},
					ParameterSchemas: commonSchema,
				}),
			),
		}),
		nil
}
