package rediscache

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			&service.ServiceProperties{
				ID:          "0346088a-d4b2-4478-aa32-f18e295ec1d9",
				Name:        "azure-rediscache",
				Description: "Azure Redis Cache (Experimental)",
				Bindable:    true,
				Tags:        []string{"Azure", "Redis", "Cache", "Database"},
			},
			m.serviceManager,
			service.NewPlan(&service.PlanProperties{
				ID:          "362b3d1b-5b57-4289-80ad-4a15a760c29c",
				Name:        "basic",
				Description: "Basic Tier, 250MB Cache",
				Free:        false,
				Extended: map[string]interface{}{
					"redisCacheSKU":      "Basic",
					"redisCacheFamily":   "C",
					"redisCacheCapacity": 0,
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "4af8bbd1-962d-4e26-84f1-f72d1d959d87",
				Name:        "standard",
				Description: "Standard Tier, 1GB Cache",
				Free:        false,
				Extended: map[string]interface{}{
					"redisCacheSKU":      "Standard",
					"redisCacheFamily":   "C",
					"redisCacheCapacity": 1,
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "b1057a8f-9a01-423a-bc35-e168d5c04cf0",
				Name:        "premium",
				Description: "Premium Tier, 6GB Cache",
				Free:        false,
				Extended: map[string]interface{}{
					"redisCacheSKU":      "Premium",
					"redisCacheFamily":   "P",
					"redisCacheCapacity": 1,
				},
			}),
		),
	}), nil
}
