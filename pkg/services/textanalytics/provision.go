package textanalytics

import (
	"context"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (s *serviceManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", s.preProvision),
		service.NewProvisioningStep("deployARMTemplate", s.deployARMTemplate),
	)
}

func (s *serviceManager) preProvision(
	context.Context,
	service.Instance,
) (service.InstanceDetails, error) {
	return &instanceDetails{
		ARMDeploymentName: uuid.NewV4().String(),
		TextAnalyticsName: uuid.NewV4().String(),
	}, nil
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*instanceDetails)
	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}

	outputs, err := s.armDeployer.Deploy(

		dt.ARMDeploymentName,
		instance.ProvisioningParameters.GetString("resourceGroup"),
		instance.ProvisioningParameters.GetString("location"),
		armTemplateBytes,
		map[string]interface{}{
			"location": instance.ProvisioningParameters.GetString("location"), // nolint: lll
			"name":     dt.TextAnalyticsName,
			"tier":     instance.Plan.GetProperties().Extended["textAnalyticsSku"],
		},
		map[string]interface{}{},
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}
	var ok bool

	dt.TextAnalyticsKey, ok = outputs["cognitivekey"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving key from deployment: %s",
			err,
		)
	}

	dt.Endpoint, ok = outputs["endpoint"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving endpoint from deployment: %s",
			err,
		)
	}

	dt.TextAnalyticsName, ok = outputs["name"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving endpoint from deployment: %s",
			err,
		)
	}

	return dt, err
}
