package fake

import "github.com/Azure/open-service-broker-azure/pkg/service"

const (
	// ServiceID is the service ID of the fake service
	ServiceID = "cdd1fb7a-d1e9-49e0-b195-e0bab747798a"
	// StandardPlanID is the plan ID for the standard (and only) variant of the
	// fake service
	StandardPlanID = "bd15e6f3-4ff5-477c-bb57-26313a368e74"
)

// GetCatalog returns a Catalog of service/plans offered by a module
func (m *Module) GetCatalog() (service.Catalog, error) {
	return service.NewCatalog([]service.Service{
		service.NewService(
			&service.ServiceProperties{
				ID:          ServiceID,
				Name:        "fake",
				Description: "Fake Service",
				Metadata: &service.ServiceMetadata{
					DisplayName:      "fake",
					ImageUrl:         "fake",
					LongDescription:  "Fake Service",
					DocumentationUrl: "fake",
					SupportUrl:       "fake",
				},
				Bindable: true,
				Tags:     []string{"Fake"},
			},
			m.ServiceManager,
			service.NewPlan(&service.PlanProperties{
				ID:          StandardPlanID,
				Name:        "standard",
				Description: "The ONLY sort of fake service-- one that's fake!",
				Free:        false,
			}),
		),
	}), nil
}
