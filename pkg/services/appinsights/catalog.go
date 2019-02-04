package appinsights

import "github.com/Azure/open-service-broker-azure/pkg/service"

// nolint: lll
func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			service.ServiceProperties{
				ID:          "66130ee7-451b-4c61-8b78-d5c426a06f3e",
				Name:        "azure-appinsights",
				Description: "Azure Application Insights (Experimental)",
				Metadata: service.ServiceMetadata{
					DisplayName: "Azure Application Insights",
					ImageURL: "https://raw.githubusercontent.com/MicrosoftDocs/" +
						"azure-docs/9eb1f875f3823af85e41ebc97e31c5b7202bf419/articles/media" +
						"/index/ApplicationInsights.svg?sanitize=true",
					LongDescription:  "Extensible Application Performance Management (APM) service for web developers on multiple platforms (Experimental)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/application-insights/app-insights-overview",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "Application", "Insights"},
			},
			m.serviceManager,
			service.NewPlan(service.PlanProperties{
				ID:          "c14826d6-87c4-45de-94a0-52fad0893799",
				Name:        "asp-dot-net-web",
				Description: "For ASP.NET web applications",
				Free:        false,
				Stability:   service.StabilityExperimental,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "ASP.NET web",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(),
					},
				},
				Extended: map[string]interface{}{
					"applicationType": "web",
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:          "be62781f-5750-49c1-be8f-56e9c804c5fa",
				Name:        "java-web",
				Description: "For Java web applications",
				Free:        false,
				Stability:   service.StabilityExperimental,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Java web",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(),
					},
				},
				Extended: map[string]interface{}{
					"applicationType": "java",
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:          "7dc0b4ee-2322-4fec-88d1-1cce63e47fd8",
				Name:        "node-dot-js",
				Description: "For Node.JS applications",
				Free:        false,
				Stability:   service.StabilityExperimental,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Node.JS",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(),
					},
				},
				Extended: map[string]interface{}{
					"applicationType": "Node.JS",
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:          "a75e8854-591a-4ef2-b3f1-b311d2a02902",
				Name:        "general",
				Description: "General",
				Free:        false,
				Stability:   service.StabilityExperimental,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "Default",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(),
					},
				},
				Extended: map[string]interface{}{
					"applicationType": "other",
				},
			}),
			service.NewPlan(service.PlanProperties{
				ID:          "b730dc3e-6928-4d05-9193-edff71790095",
				Name:        "app-center",
				Description: "For Mobile applications",
				Free:        false,
				Stability:   service.StabilityExperimental,
				Metadata: service.ServicePlanMetadata{
					DisplayName: "App Center",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: generateProvisioningParamsSchema(),
					},
				},
				Extended: map[string]interface{}{
					"applicationType": "MobileCenter",
				},
			}),
		),
	}), nil
}
