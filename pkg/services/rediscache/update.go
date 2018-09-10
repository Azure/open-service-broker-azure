package rediscache

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *serviceManager) ValidateUpdatingParameters(
	instance service.Instance,
) error {
	pp := instance.ProvisioningParameters
	up := instance.UpdatingParameters

	provisionCapacity := pp.GetInt64("skuCapacity")
	updateCapacity := up.GetInt64("skuCapacity")
	if provisionCapacity > updateCapacity {
		return fmt.Errorf("can not update an instance from larger capacity %d to"+
			"smaller capacity %d", provisionCapacity, updateCapacity)
	}
	return nil
}

func (s *serviceManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("updateARMTemplate", s.updateARMTemplate),
	)
}

func (s *serviceManager) updateARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*instanceDetails)
	up := instance.UpdatingParameters
	tagsObj := up.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}

	_, err := s.armDeployer.Update(
		dt.ARMDeploymentName,
		instance.ProvisioningParameters.GetString("resourceGroup"),
		instance.ProvisioningParameters.GetString("location"),
		armTemplateBytes,
		buildGoTemplate(instance, update),
		map[string]interface{}{},
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error updating redis instance %s", err)
	}

	return nil, nil
}
