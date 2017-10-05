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
					"skuSizeMB":      51200,
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
					"skuSizeMB":      51200,
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "edc2badc-d93b-4d9c-9d8e-da2f1c8c3e1c",
				Name:        "standard100",
				Description: "Standard Tier, 100 DTUs",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "MYSQLS100",
					"skuTier":        "Standard",
					"skuCapacityDTU": 100,
					"skuSizeMB":      128000,
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "9995c891-48ba-46cc-8dae-83595c1f443f",
				Name:        "standard200",
				Description: "Standard Tier, 200 DTUs",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "MYSQLS200",
					"skuTier":        "Standard",
					"skuCapacityDTU": 200,
					"skuSizeMB":      128000,
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "ae3cd3dd-9818-48c0-9cd0-62c3b130944e",
				Name:        "standard400",
				Description: "Standard Tier, 400 DTUs",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "MYSQLS400",
					"skuTier":        "Standard",
					"skuCapacityDTU": 400,
					"skuSizeMB":      128000,
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "08e4b43a-36bc-447e-a81f-8202b13e339c",
				Name:        "standard800",
				Description: "Standard Tier, 800 DTUs",
				Free:        false,
				Extended: map[string]interface{}{
					"skuName":        "MYSQLS800",
					"skuTier":        "Standard",
					"skuCapacityDTU": 800,
					"skuSizeMB":      128000,
				},
			}),
		),
	}), nil
}
