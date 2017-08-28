package mysql

import "github.com/Azure/azure-service-broker/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			&service.ServiceProperties{
				ID:          "997b8372-8dac-40ac-ae65-758b4a5075a5",
				Name:        "azure-mysqldb",
				Description: "Azure Database for MySQL Service",
				Bindable:    true,
				Tags:        []string{"Azure", "MySQL", "Database"},
			},
			service.NewPlan(&service.PlanProperties{
				ID:          "427559f1-bf2a-45d3-8844-32374a3e58aa",
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
				ID:          "1a538e06-9bcc-4077-8480-966cbf85bf36",
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
