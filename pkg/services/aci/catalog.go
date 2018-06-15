// +build experimental

package aci

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			&service.ServiceProperties{
				ID:          "451d5d19-4575-4d4a-9474-116f705ecc95",
				Name:        "azure-aci",
				Description: "Azure Container Instances (Experimental)",
				Metadata: &service.ServiceMetadata{
					DisplayName: "Azure Container Instances",
					ImageURL: "https://azure.microsoft.com/svghandler/container-instances/" +
						"?width=200",
					LongDescription: "Easily run containers on Azure with a single command" +
						" (Experimental)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/" +
						"container-instances/",
					SupportURL: "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "Container", "Instance"},
			},
			m.serviceManager,
			service.NewPlan(&service.PlanProperties{
				ID:          "d48798e2-21db-405b-abc7-aa6f0ff08f6c",
				Name:        "aci",
				Description: "Azure Container Instances",
				Free:        false,
				Stability:   service.StabilityExperimental,
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Azure Container Instances",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.serviceManager.getProvisionParametersSchema(), // nolint: lll
					},
				},
			}),
		),
	}), nil
}
