package postgresql

import "github.com/Azure/azure-service-broker/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			&service.ServiceProperties{
				ID:          "b43b4bba-5741-4d98-a10b-17dc5cee0175",
				Name:        "azure-postgresqldb",
				Description: "Azure Database for PostgreSQL Service",
				Bindable:    true,
				Tags:        []string{"Azure", "PostgreSQL", "Database"},
			},
			service.NewPlan(&service.PlanProperties{
				ID:          "b2ed210f-6a10-4593-a6c4-964e6b6fad62",
				Name:        "basic50",
				Description: "Basic Tier, 50 DTUs",
				Free:        false,
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "843d7d03-9306-447e-8c19-25ccc4ac30d7",
				Name:        "basic100",
				Description: "Basic Tier, 100 DTUs",
				Free:        false,
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "587e18b8-db4c-4496-87f6-2da97338aa9e",
				Name:        "standard100",
				Description: "Standard Tier, 100 DTUs",
				Free:        false,
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "185167a4-0f0c-4979-8473-6ba2f16f249d",
				Name:        "standard200",
				Description: "Standard Tier, 200 DTUs",
				Free:        false,
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "7040aa7f-ddf9-4e82-9af0-b5522969e4c7",
				Name:        "standard400",
				Description: "Standard Tier, 400 DTUs",
				Free:        false,
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "5fda2ad7-37a0-4f8b-98bf-92408c06ccfa",
				Name:        "standard800",
				Description: "Standard Tier, 800 DTUs",
				Free:        false,
			}),
		),
	}), nil
}
