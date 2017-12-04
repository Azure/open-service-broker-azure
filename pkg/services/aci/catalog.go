package aci

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			&service.ServiceProperties{
				ID:          "451d5d19-4575-4d4a-9474-116f705ecc95",
				Name:        "azure-aci",
				Description: "Azure Container Instance (Experimental)",
				Bindable:    true,
				Tags:        []string{"Azure", "Container", "Instance"},
			},
			m.serviceManager,
			service.NewPlan(&service.PlanProperties{
				ID:          "d48798e2-21db-405b-abc7-aa6f0ff08f6c",
				Name:        "aci",
				Description: "Azure Container Instances",
				Free:        false,
			}),
		),
	}), nil
}
