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
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*eventHubInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *eventHubInstanceDetails",
		)
	}
	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.EventHubName = uuid.NewV4().String()
	dt.EventHubNamespace = "eh-" + uuid.NewV4().String()
	return dt, nil
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
	plan service.Plan,
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*eventHubInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *eventHubInstanceDetails",
		)
	}
	outputs, err := s.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateBytes,
		nil, // Go template params
		map[string]interface{}{ // ARM template params
			"eventHubName":      dt.EventHubName,
			"eventHubNamespace": dt.EventHubNamespace,
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
	dt.ConnectionString = connectionString

	primaryKey, ok := outputs["primaryKey"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving primary key from deployment: %s",
			err,
		)
	}
	dt.PrimaryKey = primaryKey

	return dt, nil
}
