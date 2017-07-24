package echo

import "github.com/Azure/azure-service-broker/pkg/service"

const (
	// ServiceID is the service ID of the echo service
	ServiceID = "470b4bb6-8603-432d-aa34-d2ee74d7966c"
	// StandardPlanID is the plan ID for the standard (and only) variant of the
	// echo service
	StandardPlanID = "39ce8f26-d87d-4fb7-b06b-56f48215e308"
)

func (m *module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			&service.ServiceProperties{
				ID:          ServiceID,
				Name:        "echo",
				Description: "Echo Service",
				Bindable:    true,
				Tags:        []string{"Echo"},
			},
			service.NewPlan(&service.PlanProperties{
				ID:          StandardPlanID,
				Name:        "standard",
				Description: "The ONLY sort of echo service-- one that echoes stuff",
				Free:        false,
			}),
		),
	}), nil
}
