package aci

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (s *serviceManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
	_ service.SecureProvisioningParameters,
) error {
	pp, ok := provisioningParameters.(*ProvisioningParameters)
	if !ok {
		return errors.New(
			"error casting provisioningParameters as *aci.ProvisioningParameters",
		)
	}
	if pp.ImageName == "" {
		return service.NewValidationError(
			"image",
			fmt.Sprintf(`invalid image: "%s"`, pp.ImageName),
		)
	}
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
	dt, ok := instance.Details.(*aciInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *aciInstanceDetails",
		)
	}
	dt.ARMDeploymentName = uuid.NewV4().String()
	dt.ContainerName = uuid.NewV4().String()
	return dt, instance.SecureDetails, nil
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt, ok := instance.Details.(*aciInstanceDetails)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.Details as *aciInstanceDetails",
		)
	}
	pp, ok := instance.ProvisioningParameters.(*ProvisioningParameters)
	if !ok {
		return nil, nil, errors.New(
			"error casting instance.ProvisioningParameters as " +
				"*aci.ProvisioningParameters",
		)
	}

	outputs, err := s.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateBytes,
		pp, // Go template params
		map[string]interface{}{ // ARM template params
			"name":       dt.ContainerName,
			"image":      pp.ImageName,
			"cpuCores":   pp.NumberCores,
			"memoryInGb": fmt.Sprintf("%f", pp.Memory),
		},
		instance.Tags,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	// We don't check if this is ok, because "no public IP" is a legitimate
	// scenario.
	publicIPv4Address, ok := outputs["publicIPv4Address"].(string)
	if ok {
		dt.PublicIPv4Address = publicIPv4Address
	}

	return dt, instance.SecureDetails, nil
}
