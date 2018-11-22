package appinsights

import (
	"context"
	"fmt"
	"strings"

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
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	pp := instance.ProvisioningParameters
	var appInsightsName string
	if pp.GetString("appInsightsName") != "" {
		appInsightsName = pp.GetString("appInsightsName")
		// Check name availability
		if _, err := s.appInsightsClient.Get(
			ctx,
			pp.GetString("resourceGroup"),
			appInsightsName,
		); err != nil {
			if !strings.Contains(err.Error(), "ResourceNotFound") {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("App Insights %s already exists "+
				"in the resource group %s",
				appInsightsName,
				pp.GetString("resourceGroup"),
			)
		}
	} else {
		appInsightsName = uuid.NewV4().String()
	}

	return &instanceDetails{
		ARMDeploymentName: uuid.NewV4().String(),
		AppInsightsName:   appInsightsName,
	}, nil
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*instanceDetails)
	pp := instance.ProvisioningParameters
	tagsObj := pp.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}
	applicationType :=
		instance.Plan.GetProperties().Extended["applicationType"].(string)
	outputs, err := s.armDeployer.Deploy(
		dt.ARMDeploymentName,
		pp.GetString("resourceGroup"),
		pp.GetString("location"),
		armTemplateBytes,
		map[string]interface{}{
			"location":        pp.GetString("location"),
			"appInsightsName": dt.AppInsightsName,
			"applicationType": applicationType,
		},
		map[string]interface{}{},
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	var ok bool
	instrumentationKey, ok := outputs["instrumentationKey"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving instrumentation key from deployment: %s",
			err,
		)
	}
	dt.InstrumentationKey = service.SecureString(instrumentationKey)

	return dt, err
}
