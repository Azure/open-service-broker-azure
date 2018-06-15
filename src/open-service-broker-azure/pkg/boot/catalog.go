package boot

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

// GetCatalog returns a fully initialized catalog
func GetCatalog(
	catalogConfig service.CatalogConfig,
	azureConfig azure.Config,
) (service.Catalog, error) {
	modules, err := getModules(azureConfig)
	if err != nil {
		return nil, fmt.Errorf("error getting modules: %s", err)
	}

	// Consolidate the catalogs from all the individual modules into a single
	// catalog. Check as we go along to make sure that no two modules provide
	// services having the same ID.
	services := []service.Service{}
	usedServiceIDs := map[string]string{}
	for _, module := range modules {
		moduleName := module.GetName()
		catalog, err := module.GetCatalog()
		if err != nil {
			return nil, fmt.Errorf(
				`error retrieving catalog from module "%s": %s`,
				moduleName,
				err,
			)
		}
		for _, svc := range catalog.GetServices() {
			serviceID := svc.GetID()
			if moduleNameForUsedServiceID, ok := usedServiceIDs[serviceID]; ok {
				return nil, fmt.Errorf(
					`modules "%s" and "%s" both provide a service with the id "%s"`,
					moduleNameForUsedServiceID,
					moduleName,
					serviceID,
				)
			}

			filteredPlans := []service.Plan{}
			for _, plan := range svc.GetPlans() {
				if plan.GetStability() >= catalogConfig.MinStability {
					filteredPlans = append(filteredPlans, plan)
				}
			}
			if len(filteredPlans) > 0 {
				svc.SetPlans(filteredPlans)
				services = append(services, service.NewService(
					svc.GetProperties(),
					svc.GetServiceManager(),
					filteredPlans...,
				))
				usedServiceIDs[serviceID] = moduleName
			}
		}
	}
	catalog := service.NewCatalog(services)

	return catalog, nil
}
