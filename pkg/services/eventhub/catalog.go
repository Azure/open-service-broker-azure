package eventhub

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			&service.ServiceProperties{
				ID:          "7bade660-32f1-4fd7-b9e6-d416d975170b",
				Name:        "azure-eventhub",
				Description: "Azure Event Hub (Alpha)",
				Bindable:    true,
				Tags:        []string{"Azure", "Event", "Hubs"},
			},
			m.serviceManager,
			service.NewPlan(&service.PlanProperties{
				ID:          "80756db5-a20c-495d-ae70-62cf7d196a3c",
				Name:        "basic",
				Description: "Basic Tier, 1 Consumer group, 100 Brokered connections",
				Free:        false,
				Extended: map[string]interface{}{
					"eventHubSku": "Basic",
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:   "264ab981-9e37-44ba-b6bb-2d0fe3e80565",
				Name: "standard",
				Description: "Standard Tier, 20 Consumer groups, " +
					"1000 Brokered connections, " +
					"Additional Storage, Publisher Policies",
				Free: false,
				Extended: map[string]interface{}{
					"eventHubSku": "Standard",
				},
			}),
		),
	}), nil
}
