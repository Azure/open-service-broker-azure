package cosmosdb

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
			service.NewService(
				&service.ServiceProperties{
					ID:          "58d9fbbd-7041-4dbe-aabe-6268cd31de84",
					Name:        "azure-cosmosdb-sql",
					Description: "Azure Cosmos DB (SQL API Database Account and Database)",
					Metadata: &service.ServiceMetadata{
						DisplayName: "Azure Cosmos DB (SQL API Database Account and Database)",
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
				m.sqlAllInOneManager,
				service.NewPlan(&service.PlanProperties{
					ID:          "58d7223d-934e-4fb5-a046-0c67781eb24e",
					Name:        "sql-api",
					Description: "Azure CosmosDB With SQL API (Database Account and Database)",
					Free:        false,
					Stability:   service.StabilityExperimental,
					Metadata: &service.ServicePlanMetadata{
						DisplayName: "Azure CosmosDB (SQL API Database Account and Database)",
					},
					Schemas: service.PlanSchemas{
						ServiceInstances: service.InstanceSchemas{
							ProvisioningParametersSchema: m.sqlAllInOneManager.getProvisionParametersSchema(), // nolint: lll
						},
					},
				}),
			),
			service.NewService(
				&service.ServiceProperties{
					ID:             "6330de6f-a561-43ea-a15e-b99f44d183e6",
					Name:           "azure-cosmosdb-sql-account",
					Description:    "Azure Cosmos DB Database Account (SQL API)",
					ChildServiceID: "87c5132a-6d76-40c6-9621-0c7b7542571b",
					Metadata: &service.ServiceMetadata{
						DisplayName: "Azure Cosmos DB (SQL API - Database Account Only)",
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
				m.sqlAccountManager,
				service.NewPlan(&service.PlanProperties{
					ID:          "71168d1a-c704-49ff-8c79-214dd3d6f8eb",
					Name:        "account",
					Description: "Database Account with the SQL API",
					Free:        false,
					Stability:   service.StabilityExperimental,
					Metadata: &service.ServicePlanMetadata{
						DisplayName: "Azure CosmosDB (SQL API - Database Account Only)",
					},
					Schemas: service.PlanSchemas{
						ServiceInstances: service.InstanceSchemas{
							ProvisioningParametersSchema: m.sqlAccountManager.getProvisionParametersSchema(), // nolint: lll
						},
					},
				}),
			),
			service.NewService(
				&service.ServiceProperties{
					ID:          "87c5132a-6d76-40c6-9621-0c7b7542571b",
					Name:        "azure-cosmosdb-sql-database",
					Description: "Azure Cosmos DB Database (SQL API - Database Only)",
					Metadata: &service.ServiceMetadata{
						DisplayName: "Azure Cosmos DB (SQL API - Database Only)",
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
					ParentServiceID: "6330de6f-a561-43ea-a15e-b99f44d183e6",
				},
				m.sqlDatabaseManager,
				service.NewPlan(&service.PlanProperties{
					ID:          "c821c68c-c8e0-4176-8cf2-f0ca582a07a3",
					Name:        "database",
					Description: "Azure CosmosDB (SQL API - Database only)",
					Free:        false,
					Stability:   service.StabilityExperimental,
					Metadata: &service.ServicePlanMetadata{
						DisplayName: "Azure CosmosDB (SQL API - Database only)",
					},
				}),
			),
			service.NewService(
				&service.ServiceProperties{
					ID:          "8797a079-5346-4e84-8018-b7d5ea5c0e3a",
					Name:        "azure-cosmosdb-mongo-account",
					Description: "Azure Cosmos DB Database Account (MongoDB API)",
					Metadata: &service.ServiceMetadata{
						DisplayName: "Azure Cosmos DB (MongoDB API)",
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
					Stability:   service.StabilityExperimental,
					Metadata: &service.ServicePlanMetadata{
						DisplayName: "Azure Cosmos DB (MongoDB API)",
					},
					Schemas: service.PlanSchemas{
						ServiceInstances: service.InstanceSchemas{
							ProvisioningParametersSchema: m.mongoAccountManager.getProvisionParametersSchema(), // nolint: lll
						},
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
				},
				m.graphAccountManager,
				service.NewPlan(&service.PlanProperties{
					ID:          "126a2c47-11a3-49b1-833a-21b563de6c04",
					Name:        "account",
					Description: "Database Account with the Graph API",
					Free:        false,
					Stability:   service.StabilityExperimental,
					Metadata: &service.ServicePlanMetadata{
						DisplayName: "Azure Cosmos DB (Graph API)",
					},
					Schemas: service.PlanSchemas{
						ServiceInstances: service.InstanceSchemas{
							ProvisioningParametersSchema: m.graphAccountManager.getProvisionParametersSchema(), // nolint: lll
						},
					},
				}),
			),
			service.NewService(
				&service.ServiceProperties{
					ID:          "37915cad-5259-470d-a7aa-207ba89ada8c",
					Name:        "azure-cosmosdb-table-account",
					Description: "Azure Cosmos DB Database Account (Table API)",
					Metadata: &service.ServiceMetadata{
						DisplayName: "Azure Cosmos DB (Table API)",
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
				},
				m.graphAccountManager,
				service.NewPlan(&service.PlanProperties{
					ID:          "c970b1e8-794f-4d7c-9458-d28423c08856",
					Name:        "account",
					Description: "Database Account with the Table API",
					Free:        false,
					Stability:   service.StabilityExperimental,
					Metadata: &service.ServicePlanMetadata{
						DisplayName: "Azure Cosmos DB (Table API)",
					},
					Schemas: service.PlanSchemas{
						ServiceInstances: service.InstanceSchemas{
							ProvisioningParametersSchema: m.tableAccountManager.getProvisionParametersSchema(), // nolint: lll
						},
					},
				}),
			),
		}),
		nil
}
