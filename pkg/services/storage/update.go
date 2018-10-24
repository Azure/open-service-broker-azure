package storage

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *serviceManager) ValidateUpdatingParameters(service.Instance) error {
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
		up.GetString("resourceGroup"),
		up.GetString("location"),
		armTemplateBytes,
		buildGoTemplate(instance, *up),
		map[string]interface{}{},
		tags,
	)

	if err != nil {
		return nil, fmt.Errorf("error updating storage account %s", err)
	}
	return dt, nil
}
