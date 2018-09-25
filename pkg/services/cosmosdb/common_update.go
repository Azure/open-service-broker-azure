package cosmosdb

import (
	"context"
	"fmt"
	"reflect"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (
	c *cosmosAccountManager,
) ValidateUpdatingParameters(instance service.Instance) error {
	pp := instance.ProvisioningParameters
	up := instance.UpdatingParameters

	if err := validateReadLocations(
		"graph account update",
		up.GetStringArray("readRegions"),
	); err != nil {
		return err
	}

	// Can't update readRegions and other properties at the same time
	ppData := make(map[string]interface{})
	upData := make(map[string]interface{})
	for k, v := range pp.Data {
		ppData[k] = v
	}
	for k, v := range up.Data {
		upData[k] = v
	}
	ppData["readRegions"] = nil
	upData["readRegions"] = nil
	if !reflect.DeepEqual(
		pp.GetStringArray("readRegions"),
		up.GetStringArray("readRegions"),
	) && !reflect.DeepEqual(ppData, upData) {
		return fmt.Errorf("can't update readRegions and other properties at the same time") // nolint: lll
	}

	return nil
}

func (c *cosmosAccountManager) updateDeployment(
	pp *service.ProvisioningParameters,
	up *service.ProvisioningParameters,
	dt *cosmosdbInstanceDetails,
	kind string,
	capability string,
	additionalTags map[string]string,
) error {
	p, err := c.buildGoTemplateParams(up, dt, kind)
	if err != nil {
		return err
	}
	if capability != "" {
		p["capability"] = capability
	}
	tags := getTags(pp)
	for k, v := range additionalTags {
		tags[k] = v
	}
	err = c.deployUpdatedARMTemplate(
		up,
		dt,
		p,
		tags,
	)
	return err
}

func (c *cosmosAccountManager) deployUpdatedARMTemplate(
	pp *service.ProvisioningParameters,
	dt *cosmosdbInstanceDetails,
	goParams map[string]interface{},
	tags map[string]string,
) error {
	_, err := c.armDeployer.Update(
		dt.ARMDeploymentName,
		pp.GetString("resourceGroup"),
		pp.GetString("location"),
		armTemplateBytes,
		goParams, // Go template params
		map[string]interface{}{},
		tags,
	)
	return err
}

// This function is the same as `c.waitForReadLocationsReady` except that
// it uses `readRegions` array in updating parameter.
func (c *cosmosAccountManager) waitForReadLocationsReadyInUpdate(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*cosmosdbInstanceDetails)
	resourceGroupName := instance.ProvisioningParameters.GetString("resourceGroup")
	accountName := dt.DatabaseAccountName
	databaseAccountClient := c.databaseAccountsClient

	err := pollingUntilReadLocationsReady(
		ctx,
		resourceGroupName,
		accountName,
		databaseAccountClient,
		instance.ProvisioningParameters.GetString("location"),
		instance.UpdatingParameters.GetStringArray("readRegions"),
		false,
	)
	if err != nil {
		return nil, err
	}
	return dt, nil
}
