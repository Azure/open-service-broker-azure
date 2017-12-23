package servicebus

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
	_ service.Instance, // Reference instance
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*serviceBusInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *serviceBusInstanceDetails",
		)
	}
	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.ServiceBusNamespaceName = "sb-" + uuid.NewV4().String()
	return dt, nil
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
	plan service.Plan,
	_ service.Instance, // Reference instance
) (service.InstanceDetails, error) {
	dt, ok := instance.Details.(*serviceBusInstanceDetails)
	if !ok {
		return nil, errors.New(
			"error casting instance.Details as *serviceBusInstanceDetails",
		)
	}
	outputs, err := s.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateBytes,
		nil, // Go template params
		map[string]interface{}{ // ARM template params
			"serviceBusNamespaceName": dt.ServiceBusNamespaceName,
			"serviceBusSku":           plan.GetProperties().Extended["serviceBusSku"],
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
