package textanalytics

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			service.ServiceProperties{
				ID:          "8f6c848a-4ce1-4a69-9248-63545d3e7e9c",
				Name:        "azure-text-analytics",
				Description: "Azure Text Analytics (Experimental)",
				Metadata: service.ServiceMetadata{
					DisplayName: "Azure Text Analytics",
					ImageURL: "https://raw.githubusercontent.com/MicrosoftDocs/" +
						"azure-docs/9eb1f875f3823af85e41ebc97e31c5b7202bf419/articles/media" +
						"/index/TextAnalyticsAPI.svg?sanitize=true",
					LongDescription: "Infuse your apps, websites and bots with " +
						"intelligent algorithms to see, hear, speak, understand and " +
						"interpret your user needs through natural methods of communication." +
						" (Experimental)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/" +
						"cognitive-services/text-analytics/",
					SupportURL: "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "Cognitive", "Text Analytics", "Analytics"},
			},
			m.serviceManager,
			service.NewPlan(service.PlanProperties{
				ID:   "d5a0f91f-10da-42fc-b792-656a616d9ec2",
				Name: "free",
				Description: "Text Analytics Free Tier - max 5,000" +
					" transactions per 30 days.",
				Free: true,
				Extended: map[string]interface{}{
					"textAnalyticsSku": "F0",
				},
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Free Tier",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(),
					},
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:   "7f49713b-2689-4c66-bac9-85a024c0fb9e",
				Name: "standard-s0",
				Description: "Text Analytics Standard 0 Tier - max 25,000" +
					" transactions per 30 days.",
				Free: true,
				Extended: map[string]interface{}{
					"textAnalyticsSku": "S0",
				},
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Standard 0 Tier",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(),
					},
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:   "55575612-482b-4260-b67e-69be36d83a54",
				Name: "standard-s1",
				Description: "Text Analytics Standard 1 Tier - max 100,000" +
					" transactions per 30 days.",
				Free: true,
				Extended: map[string]interface{}{
					"textAnalyticsSku": "S1",
				},
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Standard 1 Tier",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(),
					},
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:   "76bc6a2f-1364-4ef2-8037-d7cfff48f3b6",
				Name: "standard-s2",
				Description: "Text Analytics Standard 2 Tier - max 500,000" +
					" transactions per 30 days.",
				Free: true,
				Extended: map[string]interface{}{
					"textAnalyticsSku": "S2",
				},
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Standard 2 Tier",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(),
					},
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:   "1d9a9e7c-80ac-4f23-aabe-876125541f59",
				Name: "standard-s3",
				Description: "Text Analytics Standard 3 Tier - max 2,500,000" +
					" transactions per 30 days.",
				Free: true,
				Extended: map[string]interface{}{
					"textAnalyticsSku": "S3",
				},
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Standard 3 Tier",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(),
					},
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:   "b9db834d-1350-4c50-adaf-f1e59efa2381",
				Name: "standard-s4",
				Description: "Text Analytics Standard 4 Tier - max 10,000,000" +
					" transactions per 30 days.",
				Free: true,
				Extended: map[string]interface{}{
					"textAnalyticsSku": "S4",
				},
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Standard 4 Tier",
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
