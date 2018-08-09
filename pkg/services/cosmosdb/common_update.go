package cosmosdb

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

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
	if err != nil {
		return fmt.Errorf("error deploying ARM template: %s", err)
	}
	return nil
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
	if err != nil {
		return fmt.Errorf("error deploying ARM template: %s", err)
	}
	return nil
}
