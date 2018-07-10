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
	_, _, err = c.deployARMTemplate(
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
