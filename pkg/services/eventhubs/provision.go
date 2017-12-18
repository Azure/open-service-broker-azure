package eventhubs

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
	// Nothing to validate
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
	pc, ok := instance.ProvisioningContext.(*eventHubProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting instance.ProvisioningContext as " +
				"*eventHubProvisioningContext",
		)
	}
	pc.ARMDeploymentName = uuid.NewV4().String()
	pc.EventHubName = uuid.NewV4().String()
	pc.EventHubNamespace = "eh-" + uuid.NewV4().String()
	return pc, nil
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
	plan service.Plan,
) (service.ProvisioningContext, error) {
	pc, ok := instance.ProvisioningContext.(*eventHubProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting instance.ProvisioningContext as " +
				"*eventHubProvisioningContext",
		)
	}
	outputs, err := s.armDeployer.Deploy(
		pc.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateBytes,
		nil, // Go template params
		map[string]interface{}{ // ARM template params
			"eventHubName":      pc.EventHubName,
			"eventHubNamespace": pc.EventHubNamespace,
			"eventHubSku":       plan.GetProperties().Extended["eventHubSku"],
		},
		instance.Tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	connectionString, ok := outputs["connectionString"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving connection string from deployment: %s",
			err,
		)
	}
	pc.ConnectionString = connectionString

	primaryKey, ok := outputs["primaryKey"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving primary key from deployment: %s",
			err,
		)
	}
	pc.PrimaryKey = primaryKey

	return pc, nil
}
