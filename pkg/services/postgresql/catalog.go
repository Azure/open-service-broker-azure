package postgresql

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			&service.ServiceProperties{
				ID:          "b43b4bba-5741-4d98-a10b-17dc5cee0175",
				Name:        "azure-postgresqldb",
				Description: "Azure Database for PostgreSQL (Alpha)",
				Bindable:    true,
				Tags:        []string{"Azure", "PostgreSQL", "Database"},
			},
			m.serviceManager,
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
			}),
		),
	}), nil
}
