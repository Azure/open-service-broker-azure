package cosmosdb

import (
	"context"
	"fmt"
	"strings"

	cosmosSDK "github.com/Azure/azure-sdk-for-go/services/cosmos-db/mgmt/2015-04-08/documentdb" // nolint: lll
	"github.com/Azure/open-service-broker-azure/pkg/azure/arm"
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func deleteARMDeployment(
	armDeployer arm.Deployer,
	pp *service.ProvisioningParameters,
	dt *cosmosdbInstanceDetails,
) error {
	if err := armDeployer.Delete(
		dt.ARMDeploymentName,
		pp.GetString("resourceGroup"),
	); err != nil {
		return fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return nil
}

func deleteCosmosDBAccount(
	ctx context.Context,
	databaseAccountsClient cosmosSDK.DatabaseAccountsClient,
	pp *service.ProvisioningParameters,
	dt *cosmosdbInstanceDetails,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	result, err := databaseAccountsClient.Delete(
		ctx,
		pp.GetString("resourceGroup"),
		dt.DatabaseAccountName,
	)
	if err != nil {
		return fmt.Errorf("error deleting cosmosdb server: %s", err)
	}
	if err := result.WaitForCompletion(
		ctx,
		databaseAccountsClient.Client,
	); err != nil {
		// Workaround for https://github.com/Azure/azure-sdk-for-go/issues/759
		if !strings.Contains(err.Error(), "StatusCode=404") {
			return fmt.Errorf("error deleting cosmosdb server: %s", err)
		}
	}
	return nil
}

func (c *cosmosAccountManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", c.deleteARMDeployment),
		service.NewDeprovisioningStep(
			"deleteCosmosDBServer",
			c.deleteCosmosDBAccount,
		),
	)
}

func (c *cosmosAccountManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	if err := deleteARMDeployment(
		c.armDeployer,
		instance.ProvisioningParameters,
		instance.Details.(*cosmosdbInstanceDetails),
	); err != nil {
		return nil, err
	}
	return instance.Details, nil
}

func (c *cosmosAccountManager) deleteCosmosDBAccount(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	if err := deleteCosmosDBAccount(
		ctx,
		c.databaseAccountsClient,
		instance.ProvisioningParameters,
		instance.Details.(*cosmosdbInstanceDetails),
	); err != nil {
		return nil, err
	}
	return instance.Details, nil
}
