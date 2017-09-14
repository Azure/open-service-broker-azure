package cosmosdb

import "github.com/Azure/azure-service-broker/pkg/service"

const kindKey = "kind"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
			service.NewService(
				&service.ServiceProperties{
					ID:   "6330de6f-a561-43ea-a15e-b99f44d183e6",
					Name: "azure-document-db",
					Description: "Azure DocumentDB provided by CosmosDB and accessible " +
						"via SQL (DocumentDB), Gremlin (Graph), and Table (Key-Value) APIs",
					Bindable: true,
					Tags: []string{"Azure",
						"CosmosDB",
						"Database",
						"SQL",
						"DocumentDB",
						"Gremlin",
						"Graph",
						"Table",
						"Key-Value",
					},
				},
				service.NewPlan(&service.PlanProperties{
					ID:   "71168d1a-c704-49ff-8c79-214dd3d6f8eb",
					Name: "DocumentDB",
					Description: "Azure DocumentDB provided by CosmosDB and accessible " +
						"via SQL (DocumentDB), Gremlin (Graph), and Table (Key-Value) APIs",
					Free: false,
					Extended: map[string]interface{}{
						kindKey: databaseKindGlobalDocumentDB,
					},
				}),
			),
			service.NewService(
				&service.ServiceProperties{
					ID:          "8797a079-5346-4e84-8018-b7d5ea5c0e3a",
					Name:        "azure-mongo-db",
					Description: "MongoDB on Azure provided by CosmosDB",
					Bindable:    true,
					Tags: []string{"Azure",
						"CosmosDB",
						"Database",
						"MongoDB",
					},
				},
				service.NewPlan(&service.PlanProperties{
					ID:          "86fdda05-78d7-4026-a443-1325928e7b02",
					Name:        "MongoDB",
					Description: "MongoDB",
					Free:        false,
					Extended: map[string]interface{}{
						kindKey: databaseKindMongoDB,
					},
				}),
			),
		}),
		nil
}
