package search

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			&service.ServiceProperties{
				ID:          "c54902aa-3027-4c5c-8e96-5b3d3b452f7f",
				Name:        "azure-search",
				Description: "Azure Search (Experimental)",
				Metadata: &service.ServiceMetadata{
					DisplayName: "Azure Search",
					ImageURL: "https://azure.microsoft.com/svghandler/search/" +
						"?width=200",
					LongDescription: "Cloud search service for web and mobile app " +
						"development (Experimental)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/search/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "Search", "Elasticsearch"},
			},
			m.serviceManager,
			service.NewPlan(&service.PlanProperties{
				ID:          "35bd6e80-5ff5-487e-be0e-338aee9321e4",
				Name:        "free",
				Description: "Free Tier. Max 3 Indexes, 50MB Storage/Partition",
				Free:        true,
				Extended: map[string]interface{}{
					"searchServiceSku": "free",
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Free Tier",
					Bullets: []string{
						"Max 3 Indexes",
						"50MB Storage/Partition",
					},
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
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Basic Tier",
					Bullets: []string{
						"Max 5 Indexes",
						"2GB Storage/Partition",
					},
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
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "S1 Tier",
					Bullets: []string{
						"Max 50 Indexes",
						"25GB Storage/Partition",
					},
				},
			}),
		),
	}), nil
}
