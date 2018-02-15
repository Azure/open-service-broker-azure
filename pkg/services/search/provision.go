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
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*searchInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *searchInstanceDetails",
		)
	}
	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.ServiceName = uuid.NewV4().String()
	return dt, instance.SecureDetails, nil
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*searchInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *searchInstanceDetails",
		)
	}
	sdt, ok := instance.SecureDetails.(*searchSecureInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.SecureDetails as *searchSecureInstanceDetails",
		)
	}
	outputs, err := s.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateBytes,
		nil, // Go template params
		map[string]interface{}{ // ARM template params
			"searchServiceName": dt.ServiceName,
			"searchServiceSku": instance.Plan.
				GetProperties().Extended["searchServiceSku"],
		},
		instance.Tags,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	serviceName, ok := outputs["searchServiceName"].(string)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error retrieving service name from deployment: %s",
			err,
		)
	}
	dt.ServiceName = serviceName

	apiKey, ok := outputs["apiKey"].(string)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error retrieving api key from deployment: %s",
			err,
		)
	}
	sdt.APIKey = apiKey

	return dt, sdt, nil
}
