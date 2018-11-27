package storage

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *storageManager) ValidateUpdatingParameters(
	instance service.Instance,
) error {
	pp := instance.ProvisioningParameters
	up := instance.UpdatingParameters
	previousAccountType := pp.GetString("accountType")
	nowAccountType := up.GetString("accountType")
	if previousAccountType != nowAccountType {
		if previousAccountType == "Standard_ZRS" {
			return fmt.Errorf("account type using ZRS can't be changed")
		} else if previousAccountType == "Premium_LRS" {
			return fmt.Errorf("account type using Premium_LRS can't be changed")
		}
	}
	return nil
}

func (s *storageManager) GetUpdater(service.Plan) (service.Updater, error) {
	return service.NewUpdater(
		service.NewUpdatingStep("updateARMTemplate", s.updateARMTemplate),
	)
}

func (s *storageManager) updateARMTemplate(
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
