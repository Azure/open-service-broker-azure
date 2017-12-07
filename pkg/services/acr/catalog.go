package acr

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			&service.ServiceProperties{
				ID:          "0b9401d6-4c04-4b10-a4da-0fd6cd1c7b4a",
				Name:        "acr",
				Description: "Azure Container Registry (Experimental)",
				Bindable:    true,
				Tags:        []string{"Azure", "Container Registry"},
			},
			m.serviceManager,
			service.NewPlan(&service.PlanProperties{
				ID:          "e74c0b8e-23b1-4495-91ae-41682d3a0b7c",
				Name:        "Basic",
				Description: "Basic Tier",
				Free:        true,
				Extended: map[string]interface{}{
					"registrySku": "Basic",
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "e642e66d-6c40-4f3b-a9f5-c2acd08eae50",
				Name:        "Standard",
				Description: "Standard Tier.",
				Free:        true,
				Extended: map[string]interface{}{
					"registrySku": "Standard",
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "589e0f3a-7095-4e4a-be37-904f291fdb54",
				Name:        "Premium",
				Description: "Premium Tier.",
				Free:        true,
				Extended: map[string]interface{}{
					"registrySku": "Premium",
				},
			}),
		),
	}), nil
}
