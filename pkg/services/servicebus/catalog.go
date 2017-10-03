package servicebus

import "github.com/Azure/azure-service-broker/pkg/service"

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			&service.ServiceProperties{
				ID:          "6dc44338-2f13-4bc5-9247-5b1b3c5462d3",
				Name:        "azure-servicebus ",
				Description: "Azure Service Bus",
				Bindable:    true,
				Tags:        []string{"Azure", "Service", "Bus"},
			},
			service.NewPlan(&service.PlanProperties{
				ID:          "13c6da8f-128c-48c0-a3a9-659d1b6d3920",
				Name:        "basic",
				Description: "Basic Tier, Shared Capacity",
				Free:        false,
				Extended: map[string]interface{}{
					"serviceBusSku": "Basic",
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:   "6be0d8b5-381f-4d68-bdfd-a131425d3835",
				Name: "standard",
				Description: "Standard Tier, Shared Capacity, Topics, 12.5M " +
					"Messaging Operations/Month, Variable Pricing",
				Free: false,
				Extended: map[string]interface{}{
					"serviceBusSku": "Standard",
				},
			}),
			service.NewPlan(&service.PlanProperties{
				ID:   "e359cbbe-d52b-47cc-8243-5bb9651c86c7",
				Name: "premium",
				Description: "Premium Tier, Dedicated Capacity, Recommended " +
					"For Production Workloads, Fixed Pricing",
				Free: false,
				Extended: map[string]interface{}{
					"serviceBusSku": "Premium",
				},
			}),
		),
	}), nil
}
