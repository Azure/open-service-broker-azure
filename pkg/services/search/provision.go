package search

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (s *serviceManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	// No validation needed
	return nil
}

func (s *serviceManager) GetProvisioner(
	service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", s.preProvision),
		service.NewProvisioningStep("deployARMTemplate", s.deployARMTemplate),
	)
}

func (s *serviceManager) preProvision(
	_ context.Context,
	instance service.Instance,
	_ service.Plan,
) (service.ProvisioningContext, error) {
	pc, ok := instance.ProvisioningContext.(*searchProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting instance.ProvisioningContext as " +
				"*searchProvisioningContext",
		)
	}
	pc.ARMDeploymentName = uuid.NewV4().String()
	pc.ServiceName = uuid.NewV4().String()
	return pc, nil
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
	plan service.Plan,
) (service.ProvisioningContext, error) {
	pc, ok := instance.ProvisioningContext.(*searchProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting instance.ProvisioningContext as " +
				"*searchProvisioningContext",
		)
	}
	outputs, err := s.armDeployer.Deploy(
		pc.ARMDeploymentName,
		instance.StandardProvisioningContext.ResourceGroup,
		instance.StandardProvisioningContext.Location,
		armTemplateBytes,
		nil, // Go template params
		map[string]interface{}{ // ARM template params
			"searchServiceName": pc.ServiceName,
			"searchServiceSku":  plan.GetProperties().Extended["searchServiceSku"],
		},
		instance.StandardProvisioningContext.Tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	serviceName, ok := outputs["searchServiceName"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving service name from deployment: %s",
			err,
		)
	}
	pc.ServiceName = serviceName

	apiKey, ok := outputs["apiKey"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving api key from deployment: %s",
			err,
		)
	}
	pc.APIKey = apiKey

	return pc, nil
}
