package keyvault

import "github.com/Azure/azure-service-broker/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			&service.ServiceProperties{
				ID:          "d90c881e-c9bb-4e07-a87b-fcfe87e03276",
				Name:        "azure-keyvault",
				Description: "Azure Key Vault",
				Bindable:    true,
				Tags:        []string{"Azure", "Key", "Vault"},
			},
			service.NewPlan(&service.PlanProperties{
				ID:          "3577ee4a-75fc-44b3-b354-9d33d52ef486",
				Name:        "standard",
				Description: "Standard Tier",
				Free:        false,
				Extended: map[string]interface{}{
					"vaultSku": "Standard",
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:          "6893b1de-0a7b-42bb-b28d-1636c4b81f75",
				Name:        "premium",
				Description: "Premium Tier",
				Free:        false,
				Extended: map[string]interface{}{
					"vaultSku": "Premium",
				},
			}),
		),
	}), nil
}
