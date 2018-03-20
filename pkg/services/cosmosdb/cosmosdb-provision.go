package cosmosdb

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (c *cosmosAccountManager) ValidateProvisioningParameters(
	service.ProvisioningParameters,
	service.SecureProvisioningParameters,
) error {
	// Nothing to validate
	return nil
}

func (c *cosmosAccountManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	return preProvision(instance)
}

func (c *cosmosAccountManager) buildGoTemplateParams(
	dt *cosmosdbInstanceDetails,
) map[string]interface{} {
	p := map[string]interface{}{}
	p["name"] = dt.DatabaseAccountName
	p["kind"] = "GlobalDocumentDB"
	return p
}

func (c *cosmosAccountManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
	goParams map[string]interface{},
) (*cosmosdbInstanceDetails, *cosmosdbSecureInstanceDetails, error) {
	dt := cosmosdbInstanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}

	outputs, err := c.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateBytes,
		goParams, // Go template params
		map[string]interface{}{},
		instance.Tags,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	var ok bool
	dt.FullyQualifiedDomainName, ok = outputs["fullyQualifiedDomainName"].(string)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error retrieving fully qualified domain name from deployment: %s",
			err,
		)
	}

	primaryKey, ok := outputs["primaryKey"].(string)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error retrieving primary key from deployment: %s",
			err,
		)
	}

	sdt := cosmosdbSecureInstanceDetails{
		PrimaryKey: primaryKey,
	}

	return &dt, &sdt, nil
}
