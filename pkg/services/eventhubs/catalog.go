package eventhubs

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			service.ServiceProperties{
				ID:          "7bade660-32f1-4fd7-b9e6-d416d975170b",
				Name:        "azure-eventhubs",
				Description: "Azure Event Hubs (Experimental)",
				Metadata: service.ServiceMetadata{
					DisplayName: "Azure Event Hubs",
					ImageURL: "https://raw.githubusercontent.com/MicrosoftDocs/" +
						"azure-docs/9eb1f875f3823af85e41ebc97e31c5b7202bf419/articles/media" +
						"/index/EventHubs.svg?sanitize=true",
					LongDescription: "Cloud-scale telemetry ingestion from websites, apps, " +
						"and any streams of data (Experimental)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/event-hubs/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "Event", "Hubs"},
			},
			m.serviceManager,
			service.NewPlan(service.PlanProperties{
				ID:          "80756db5-a20c-495d-ae70-62cf7d196a3c",
				Name:        "basic",
				Description: "Basic Tier, 1 Consumer group, 100 Brokered connections",
				Free:        false,
				Stability:   service.StabilityExperimental,
				Extended: map[string]interface{}{
					"eventHubSku": "Basic",
				},
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Basic Tier",
					Bullets: []string{
						"1 Consumer group",
						"100 Brokered connections",
					},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(),
					},
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:   "264ab981-9e37-44ba-b6bb-2d0fe3e80565",
				Name: "standard",
				Description: "Standard Tier, 20 Consumer groups, " +
					"1000 Brokered connections, " +
					"Additional Storage, Publisher Policies",
				Free:      false,
				Stability: service.StabilityExperimental,
				Extended: map[string]interface{}{
					"eventHubSku": "Standard",
				},
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Standard Tier",
					Bullets: []string{
						"20 Consumer groups",
						"1000 Brokered connections",
						"Additional Storage",
						"Publisher Policies",
					},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(),
					},
				},
			}),
		),
	}), nil
}
