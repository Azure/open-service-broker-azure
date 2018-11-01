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
					DisplayName:      "Azure Application Insights",
					ImageURL:         "https://azure.microsoft.com/svghandler/application-insights/?width=200",
					LongDescription:  "Extensible Application Performance Management (APM) service for web developers on multiple platforms (Experimental)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/application-insights/app-insights-overview",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "Application", "Insights"},
			},
			m.serviceManager,
			service.NewPlan(service.PlanProperties{
				ID:          "a75e8854-591a-4ef2-b3f1-b311d2a02902",
				Name:        "default",
				Description: "Default",
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
			}),
		),
	}), nil
}
