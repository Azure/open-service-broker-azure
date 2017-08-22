package mysql

import "github.com/Azure/azure-service-broker/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			&service.ServiceProperties{
				ID:          "e40b3635-01bc-4262-b2c5-0847bd7ab43b",
				Name:        "azure-mysqldb",
				Description: "Azure Database for MySQL Service",
				Bindable:    true,
				Tags:        []string{"Azure", "MySQL", "Database"},
			},
			service.NewPlan(&service.PlanProperties{
				ID:          "a05c9967-9f20-4b17-8c66-bb32ab396fcd",
				Name:        "basic50",
				Description: "Basic Tier, 50 DTUs.",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "MYSQLB50",
					"skuTier":        "Basic",
					"skuCapacityDTU": 50,
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "d8d5cac9-d975-48ea-9ac4-8232f92bcb93",
				Name:        "basic100",
				Description: "Basic Tier, 100 DTUs",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "MYSQLB100",
					"skuTier":        "Basic",
					"skuCapacityDTU": 100,
				},
			}),
		),
	}), nil
}
