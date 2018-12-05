package servicebus

import "github.com/Azure/open-service-broker-azure/pkg/service"

// nolint: lll
func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			service.ServiceProperties{
				ID:          "6dc44338-2f13-4bc5-9247-5b1b3c5462d3",
				Name:        "azure-servicebus-namespace",
				Description: "Azure Service Bus Namespace(Experimental)",
				// It has two childs, should be changed after refactoring.
				ChildServiceID: "0e93fbb8-7904-43a5-82db-81c7d3886a24",
				Metadata: service.ServiceMetadata{
					DisplayName: "Azure Service Bus Namespace",
					ImageURL: "https://azure.microsoft.com/svghandler/service-bus/" +
						"?width=200",
					LongDescription: "Reliable cloud messaging as a service (MaaS) and " +
						"simple hybrid integration. Create an Azure Service Bus" +
						"Namespace. (Experimental)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/service-bus/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "Service", "Bus"},
			},
			m.namespaceManager,
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
						ProvisioningParametersSchema: generateNamespaceProvisioningParamsSchema(),
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
						ProvisioningParametersSchema: generateNamespaceProvisioningParamsSchema(),
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
						ProvisioningParametersSchema: generateNamespaceProvisioningParamsSchema(),
					},
				},
			}),
		),
		service.NewService(
			service.ServiceProperties{
				ID:              "0e93fbb8-7904-43a5-82db-81c7d3886a24",
				Name:            "azure-servicebus-queue",
				Description:     "Azure Service Bus Queue(Experimental)",
				ParentServiceID: "6dc44338-2f13-4bc5-9247-5b1b3c5462d3",
				Metadata: service.ServiceMetadata{
					DisplayName: "Azure Service Bus Queue",
					ImageURL: "https://azure.microsoft.com/svghandler/service-bus/" +
						"?width=200",
					LongDescription: "Reliable cloud messaging as a service (MaaS) and " +
						"simple hybrid integration. Create an Azure Service Bus" +
						"Queue in an existing namespace. (Experimental)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/service-bus/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "Service", "Bus"},
			},
			m.queueManager,
			service.NewPlan(service.PlanProperties{
				ID:          "89440eec-a888-49ce-b392-60c653d7a98b",
				Name:        "queue",
				Description: "New queue in existing namespace",
				Free:        false,
				Stability:   service.StabilityExperimental,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Service Bus Queue",
					Bullets:     []string{"Message Queue"},
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(),
					},
				},
			}),
		),
		service.NewService(
			service.ServiceProperties{
				ID:              "dc6d1545-4391-4c7e-ac7e-a8463787fb93",
				Name:            "azure-servicebus-topic",
				Description:     "Azure Service Bus Topic(Experimental)",
				ParentServiceID: "6dc44338-2f13-4bc5-9247-5b1b3c5462d3",
				Metadata: service.ServiceMetadata{
					DisplayName: "Azure Service Bus Topic",
					ImageURL: "https://azure.microsoft.com/svghandler/service-bus/" +
						"?width=200",
					LongDescription: "Reliable cloud messaging as a service (MaaS) and " +
						"simple hybrid integration. Create an Azure Service Bus" +
						"Topic in an existing namespace. (Experimental)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/service-bus/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "Service", "Bus"},
			},
			m.topicManager,
			service.NewPlan(service.PlanProperties{
				ID:          "dd1e4d44-58be-4f34-84ff-f73ccef405e5",
				Name:        "topic",
				Description: "New topic in existing namespace",
				Free:        false,
				Stability:   service.StabilityExperimental,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Service Bus Topic",
					Bullets:     []string{"One to Many Message Mechanism"},
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
