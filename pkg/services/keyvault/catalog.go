package keyvault

import "github.com/Azure/open-service-broker-azure/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			&service.ServiceProperties{
				ID:          "d90c881e-c9bb-4e07-a87b-fcfe87e03276",
				Name:        "azure-keyvault",
				Description: "Azure Key Vault (Experimental)",
				Metadata: &service.ServiceMetadata{
					DisplayName: "Azure Key Vault",
					ImageURL: "https://azure.microsoft.com/svghandler/key-vault/" +
						"?width=200",
					LongDescription: "Safeguard cryptographic keys and other secrets " +
						"used by cloud apps and services (Experimental)",
					DocumentationURL: "https://docs.microsoft.com/en-us/azure/key-vault/",
					SupportURL:       "https://azure.microsoft.com/en-us/support/",
				},
				Bindable: true,
				Tags:     []string{"Azure", "Key", "Vault"},
			},
			m.serviceManager,
			service.NewPlan(&service.PlanProperties{
				ID:          "3577ee4a-75fc-44b3-b354-9d33d52ef486",
				Name:        "standard",
				Description: "Standard Tier",
				Free:        false,
				Stability:   service.StabilityExperimental,
				Extended: map[string]interface{}{
					"vaultSku": "Standard",
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Standard Tier",
				},
				Schemas: service.PlanSchemas{
					ServiceInstances: service.InstanceSchemas{
						ProvisioningParametersSchema: m.serviceManager.getProvisionParametersSchema(), // nolint: lll
					},
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "6893b1de-0a7b-42bb-b28d-1636c4b81f75",
				Name:        "premium",
				Description: "Premium Tier",
				Free:        false,
				Stability:   service.StabilityExperimental,
				Extended: map[string]interface{}{
					"vaultSku": "Premium",
				},
				Metadata: &service.ServicePlanMetadata{
					DisplayName: "Premium Tier",
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
