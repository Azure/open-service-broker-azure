package main

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func getCatalog(modules []service.Module) (service.Catalog, error) {
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
			services = append(services, svc)
			usedServiceIDs[serviceID] = moduleName
		}
	}
	catalog := service.NewCatalog(services)

	return catalog, nil
}
