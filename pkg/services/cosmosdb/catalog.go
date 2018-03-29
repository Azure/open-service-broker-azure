package cosmosdb

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
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
					ProvisionParamsSchema: m.sqlAccountManager.getProvisionParametersSchema(), // nolint: lll
				},
				m.sqlAccountManager,
				service.NewPlan(&service.PlanProperties{
					ID:          "71168d1a-c704-49ff-8c79-214dd3d6f8eb",
					Name:        "account",
					Description: "Database Account with the SQL API",
					Free:        false,
					Metadata: &service.ServicePlanMetadata{
						DisplayName: "Azure CosmosDB (SQL API)",
					},
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
					ProvisionParamsSchema: m.mongoAccountManager.getProvisionParametersSchema(), // nolint: lll
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
				}),
			),
			service.NewService(
				&service.ServiceProperties{
					ID:          "5f5252a0-6922-4a0c-a755-f9be70d7c79b",
					Name:        "azure-cosmosdb-graph-account",
					Description: "Azure Cosmos DB Database Account (Graph API)",
					Metadata: &service.ServiceMetadata{
						DisplayName: "Azure Cosmos DB (Graph API)",
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
						"Graph",
						"Gremlin",
					},
					ProvisionParamsSchema: m.graphAccountManager.getProvisionParametersSchema(), // nolint: lll
				},
				m.graphAccountManager,
				service.NewPlan(&service.PlanProperties{
					ID:          "126a2c47-11a3-49b1-833a-21b563de6c04",
					Name:        "account",
					Description: "Database Account with the Graph API",
					Free:        false,
					Metadata: &service.ServicePlanMetadata{
						DisplayName: "Azure CosmosDB with Graph API",
					},
				}),
			),
			service.NewService(
				&service.ServiceProperties{
					ID:          "37915cad-5259-470d-a7aa-207ba89ada8c",
					Name:        "azure-cosmosdb-table-account",
					Description: "Azure Cosmos DB Database Account (Table API)",
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
						"Table",
					},
					ProvisionParamsSchema: m.tableAccountManager.getProvisionParametersSchema(), // nolint: lll
				},
				m.graphAccountManager,
				service.NewPlan(&service.PlanProperties{
					ID:          "c970b1e8-794f-4d7c-9458-d28423c08856",
					Name:        "account",
					Description: "Database Account with the Table API",
					Free:        false,
					Metadata: &service.ServicePlanMetadata{
						DisplayName: "Azure CosmosDB with Table API",
					},
				}),
			),
		}),
		nil
}
