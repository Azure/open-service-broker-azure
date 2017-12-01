package search

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/service"
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
	_ string, // instanceID
	_ service.Plan,
	_ service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	_ service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*searchProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *searchProvisioningContext",
		)
	}
	pc.ARMDeploymentName = uuid.NewV4().String()
	pc.ServiceName = uuid.NewV4().String()
	return pc, nil
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	_ string, // instanceID
	plan service.Plan,
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	_ service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*searchProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *searchProvisioningContext",
		)
	}
	outputs, err := s.armDeployer.Deploy(
		pc.ARMDeploymentName,
		standardProvisioningContext.ResourceGroup,
		standardProvisioningContext.Location,
		armTemplateBytes,
		nil, // Go template params
		map[string]interface{}{ // ARM template params
			"searchServiceName": pc.ServiceName,
			"searchServiceSku":  plan.GetProperties().Extended["searchServiceSku"],
		},
		standardProvisioningContext.Tags,
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
