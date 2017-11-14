package search

import "github.com/Azure/azure-service-broker/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			&service.ServiceProperties{
				ID:          "c54902aa-3027-4c5c-8e96-5b3d3b452f7f",
				Name:        "azuresearch",
				Description: "Azure Search (Alpha)",
				Bindable:    true,
				Tags:        []string{"Azure", "Search", "Elasticsearch"},
			},
			service.NewPlan(&service.PlanProperties{
				ID:          "35bd6e80-5ff5-487e-be0e-338aee9321e4",
				Name:        "free",
				Description: "Free Tier. Max 3 Indexes, 50MB Storage/Partition",
				Free:        true,
				Extended: map[string]interface{}{
					"searchServiceSku": "free",
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "4a50e008-5513-42d3-8b2f-d8b3ad43f7eb",
				Name:        "basic",
				Description: "Basic Tier. Max 5 Indexes, 2GB Storage/Partition",
				Free:        true,
				Extended: map[string]interface{}{
					"searchServiceSku": "basic",
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "65e89af2-8da2-4559-b103-8dd2dd8fdd40",
				Name:        "standard-s1",
				Description: "S1 Tier. Max 50 Indexes, 25GB Storage/Partition",
				Free:        true,
				Extended: map[string]interface{}{
					"searchServiceSku": "standard",
				},
			}),
		),
	}), nil
}
