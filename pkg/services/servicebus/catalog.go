package servicebus

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			service.ServiceProperties{
				ID:          "6dc44338-2f13-4bc5-9247-5b1b3c5462d3",
				Name:        "azure-servicebus",
				Description: "Azure Service Bus (Experimental)",
				Metadata: service.ServiceMetadata{
					DisplayName: "Azure Service Bus",
					ImageURL: "https://azure.microsoft.com/svghandler/service-bus/" +
						"?width=200",
					LongDescription: "Reliable cloud messaging as a service (MaaS) and " +
						"simple hybrid integration (Experimental)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/service-bus/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "Service", "Bus"},
			},
			m.serviceManager,
			service.NewPlan(service.PlanProperties{
				ID:          "d06817b1-87ea-4320-8942-14b1d060206a",
				Name:        "basic",
				Description: "Basic Tier, Shared Capacity",
				Free:        false,
				Stability:   service.StabilityExperimental,
				Extended: map[string]interface{}{
					"serviceBusSku": "Basic",
				},
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Basic Tier",
					Bullets:     []string{"Shared Capacity"},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(),
					},
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:   "6be0d8b5-381f-4d68-bdfd-a131425d3835",
				Name: "standard",
				Description: "Standard Tier, Shared Capacity, Topics, 12.5M " +
					"Messaging Operations/Month, Variable Pricing",
				Free:      false,
				Stability: service.StabilityExperimental,
				Extended: map[string]interface{}{
					"serviceBusSku": "Standard",
				},
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Standard Tier",
					Bullets: []string{
						"Shared Capacity",
						"Topics",
						"12.5M Messaging Operations/Month",
						"Variable Pricing",
					},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(),
					},
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:   "cec378a7-6452-4203-beca-d34898edbadc",
				Name: "premium",
				Description: "Premium Tier, Dedicated Capacity, Recommended " +
					"For Production Workloads, Fixed Pricing",
				Free:      false,
				Stability: service.StabilityExperimental,
				Extended: map[string]interface{}{
					"serviceBusSku": "Premium",
				},
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Premium Tier",
					Bullets: []string{
						"Dedicated Capacity",
						"Recommended For Production Workloads",
						"Fixed Pricing",
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
