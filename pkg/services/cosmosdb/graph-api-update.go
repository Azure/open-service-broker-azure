package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	g *graphAccountManager,
) ValidateUpdatingParameters(instance service.Instance) error {
	return validateReadLocations(
		"graph account update",
		instance.UpdatingParameters.GetStringArray("readRegions"),
	)
}

func (
	g *graphAccountManager,
) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		// The cosmosDB has a contraint: it cannot update properties and
		// add/remove regions at the same time, so we must deal with the update twice,
		// one time updating region, one time updating properties.
		service.NewUpdatingStep("updateReadLocations", g.updateReadLocations),
		service.NewUpdatingStep("waitForReadLocationsReady", g.waitForReadLocationsReady), //nolint: lll
		service.NewUpdatingStep("updateARMTemplate", g.updateARMTemplate),
	)
}

func (g *graphAccountManager) updateReadLocations(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	if err := g.cosmosAccountManager.updateReadLocations(
		instance.ProvisioningParameters,
		instance.UpdatingParameters,
		instance.Details.(*cosmosdbInstanceDetails),
		"GlobalDocumentDB",
		"EnableGremlin",
		map[string]string{
			"defaultExperience": "Graph",
		},
	); err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}

func (g *graphAccountManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	if err := g.cosmosAccountManager.updateDeployment(
		instance.ProvisioningParameters,
		instance.UpdatingParameters,
		instance.Details.(*cosmosdbInstanceDetails),
		"GlobalDocumentDB",
		"EnableGremlin",
		map[string]string{
			"defaultExperience": "Graph",
		},
	); err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	return instance.Details, nil
}
