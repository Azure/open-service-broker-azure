package postgresql

import "github.com/Azure/open-service-broker-azure/pkg/service"

// nolint: lll
func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		// all-in-one
		service.NewService(
			&service.ServiceProperties{
				ID:          "b43b4bba-5741-4d98-a10b-17dc5cee0175",
				Name:        "azure-postgresql",
				Description: "Azure Database for PostgreSQL-- DBMS and single database (preview)",
				Metadata: &service.ServiceMetadata{
					DisplayName:      "Azure Database for PostgreSQL",
					ImageURL:         "https://azure.microsoft.com/svghandler/postgresql/?width=200",
					LongDescription:  "Azure Database for PostgreSQL-- DBMS and single database (preview)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/postgresql/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "PostgreSQL", "DBMS", "Server", "Database"},
				ProvisionParamsSchema: m.allInOneManager.getProvisionParametersSchema(),
			},
			m.allInOneManager,
			service.NewPlan(&service.PlanProperties{
				ID:          "b2ed210f-6a10-4593-a6c4-964e6b6fad62",
				Name:        "basic50",
				Description: "Basic Tier, 50 DTUs",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "PGSQLB50",
					"skuTier":        "Basic",
					"skuCapacityDTU": 50,
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Basic Tier",
					Bullets:     []string{"50 DTUs"},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "843d7d03-9306-447e-8c19-25ccc4ac30d7",
				Name:        "basic100",
				Description: "Basic Tier, 100 DTUs",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "PGSQLB100",
					"skuTier":        "Basic",
					"skuCapacityDTU": 100,
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Basic Tier",
					Bullets:     []string{"100 DTUs"},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "bd588e32-0514-4421-8ef3-f54039914e61",
				Name:        "standard100",
				Description: "Standard Tier, 100 DTUs",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "PGSQLS100",
					"skuTier":        "Standard",
					"skuCapacityDTU": 100,
				},
			}),
		),
		// dbms only
		service.NewService(
			&service.ServiceProperties{
				ID:             "d3f74b44-79bc-4d1e-bf7d-c247c2b851f9",
				Name:           "azure-postgresql-dbms",
				Description:    "Azure Database for PostgreSQL-- DBMS only (preview)",
				ChildServiceID: "25434f16-d762-41c7-bbdd-8045d7f74ca",
				Metadata: &service.ServiceMetadata{
					DisplayName:      "Azure Database for PostgreSQL-- DBMS Only",
					ImageURL:         "https://azure.microsoft.com/svghandler/postgresql/?width=200",
					LongDescription:  "Azure Database for PostgreSQL-- DBMS only (preview)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/postgresql/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: false,
				Tags:     []string{"Azure", "PostgreSQL", "DBMS", "Server", "Database"},
				ProvisionParamsSchema: m.dbmsManager.getProvisionParametersSchema(),
			},
			m.dbmsManager,
			service.NewPlan(&service.PlanProperties{
				ID:          "bf389028-8dcc-433a-ab6f-0ee9b8db142f",
				Name:        "basic50",
				Description: "Basic Tier, 50 DTUs",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "PGSQLB50",
					"skuTier":        "Basic",
					"skuCapacityDTU": 50,
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Basic Tier",
					Bullets:     []string{"50 DTUs"},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "58633c61-942c-42cb-b22c-346a4c594b8e",
				Name:        "basic100",
				Description: "Basic Tier, 100 DTUs",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "PGSQLB100",
					"skuTier":        "Basic",
					"skuCapacityDTU": 100,
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Basic Tier",
					Bullets:     []string{"100 DTUs"},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "ba653f35-2152-4d76-8641-b21d4636b2e1",
				Name:        "standard100",
				Description: "Standard Tier, 100 DTUs",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "PGSQLS100",
					"skuTier":        "Standard",
					"skuCapacityDTU": 100,
				},
			}),
		),
		// database only
		service.NewService(
			&service.ServiceProperties{
				ID:              "25434f16-d762-41c7-bbdd-8045d7f74ca6",
				Name:            "azure-postgresql-database",
				Description:     "Azure Database for PostgreSQL-- database only (preview)",
				ParentServiceID: "d3f74b44-79bc-4d1e-bf7d-c247c2b851f9",
				Metadata: &service.ServiceMetadata{
					DisplayName:      "Azure Database for PostgreSQL-- Database Only",
					ImageURL:         "https://azure.microsoft.com/svghandler/postgresql/?width=200",
					LongDescription:  "Azure Database for PostgreSQL-- database only (preview)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/postgresql/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "PostgreSQL", "Database"},
				ProvisionParamsSchema: m.databaseManager.getProvisionParametersSchema(),
			},
			m.databaseManager,
			service.NewPlan(&service.PlanProperties{
				ID:          "df6f5ef1-e602-406b-ba73-09c107d1e31b",
				Name:        "database",
				Description: "A new database added to an existing DBMS (preview)",
				Free:        false,
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Azure Database for PostgreSQL-- Database Only",
				},
			}),
		),
	}), nil
}
