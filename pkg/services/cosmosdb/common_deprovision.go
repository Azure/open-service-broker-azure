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
	instance service.Instance,
) error {
	dt, ok := instance.Details.(*cosmosdbInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Details as *cosmosdbInstanceDetails",
		)
	}
	if err := armDeployer.Delete(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
	); err != nil {
		return fmt.Errorf("error deleting ARM deployment: %s", err)
	}
	return nil
}

func deleteCosmosDBServer(
	ctx context.Context,
	databaseAccountsClient cosmosSDK.DatabaseAccountsClient,
	instance service.Instance,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	dt, ok := instance.Details.(*cosmosdbInstanceDetails)
	if !ok {
		return fmt.Errorf(
			"error casting instance.Details as *cosmosdbInstanceDetails",
		)
	}
	result, err := databaseAccountsClient.Delete(
		ctx,
		instance.ResourceGroup,
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
		if strings.Contains(err.Error(), "StatusCode=404") {
			return nil
		}
		return fmt.Errorf("error deleting cosmosdb server: %s", err)
	}
	return nil
}
