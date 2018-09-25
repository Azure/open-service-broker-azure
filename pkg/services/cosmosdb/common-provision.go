package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

const enabled = "enabled"
const disabled = "disabled"

func generateAccountName(location string) string {
	databaseAccountName := uuid.NewV4().String()
	// CosmosDB currently limits database account names to 50 characters,
	// which includes location and a - character. Check if we will
	// exceed this and generate a shorter random identifier if needed.
	effectiveNameLength := len(location) + len(databaseAccountName)
	if effectiveNameLength > 49 {
		nameLength := 49 - len(location)
		databaseAccountName = generate.NewIdentifierOfLength(nameLength)
		logFields := log.Fields{
			"name":   databaseAccountName,
			"length": len(databaseAccountName),
		}
		log.WithFields(logFields).Debug(
			"returning fallback database account name",
		)
	}
	return databaseAccountName
}

func (c *cosmosAccountManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	l := instance.ProvisioningParameters.GetString("location")
	return &cosmosdbInstanceDetails{
		ARMDeploymentName:   uuid.NewV4().String(),
		DatabaseAccountName: generateAccountName(l),
	}, nil
}

func (c *cosmosAccountManager) buildGoTemplateParams(
	pp *service.ProvisioningParameters,
	dt *cosmosdbInstanceDetails,
	kind string,
) (map[string]interface{}, error) {
	// In Azure portal, "region" is used to indicate a place, while in REST api,
	// "location" is used to indicate a place. So we use "readRegions" when
	// communicating with users and use "readLocations" in the code.
	readLocations := pp.GetStringArray("readRegions")
	readLocations = append([]string{pp.GetString("location")}, readLocations...)
	return c.buildGoTemplateParamsCore(
		pp,
		dt,
		kind,
		readLocations,
	)
}

// The deployment will return success once the write region is created,
// ignoring the status of read regions , so we must implement detection logic
// by ourselves.
func (c *cosmosAccountManager) waitForReadLocationsReady(
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
		instance.ProvisioningParameters.GetStringArray("readRegions"),
		true,
	)
	if err != nil {
		return nil, err
	}
	return dt, nil
}

func getTags(pp *service.ProvisioningParameters) map[string]string {
	tagsObj := pp.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	return tags
}

func (c *cosmosAccountManager) deployARMTemplate(
	pp *service.ProvisioningParameters,
	dt *cosmosdbInstanceDetails,
	goParams map[string]interface{},
	tags map[string]string,
) (string, string, error) {
	outputs, err := c.armDeployer.Deploy(
		dt.ARMDeploymentName,
		pp.GetString("resourceGroup"),
		pp.GetString("location"),
		armTemplateBytes,
		goParams, // Go template params
		map[string]interface{}{},
		tags,
	)
	if err != nil {
		return "", "", fmt.Errorf("error deploying ARM template: %s", err)
	}
	fqdn, primaryKey, err := c.handleOutput(outputs)
	if err != nil {
		return "", "", fmt.Errorf("error deploying ARM template: %s", err)
	}
	return fqdn, primaryKey, nil
}

func (c *cosmosAccountManager) handleOutput(
	outputs map[string]interface{},
) (string, string, error) {

	var ok bool
	fqdn, ok := outputs["fullyQualifiedDomainName"].(string)
	if !ok {
		return "", "", fmt.Errorf(
			"error retrieving fully qualified domain name from deployment",
		)
	}

	primaryKey, ok := outputs["primaryKey"].(string)
	if !ok {
		return "", "", fmt.Errorf("error retrieving primary key from deployment")
	}
	return fqdn, primaryKey, nil
}
