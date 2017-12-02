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
) error {
	pp, ok := provisioningParameters.(*ProvisioningParameters)
	if !ok {
		return errors.New(
			"error casting provisioningParameters as " +
				"*aci.ProvisioningParameters",
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
	_ string, // instanceID
	_ service.Plan,
	_ service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	_ service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*aciProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *aciProvisioningContext",
		)
	}
	pc.ARMDeploymentName = uuid.NewV4().String()
	pc.ContainerName = uuid.NewV4().String()
	return pc, nil
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	_ string, // instanceID
	_ service.Plan,
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	provisioningParameters service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*aciProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *aciProvisioningContext",
		)
	}
	pp, ok := provisioningParameters.(*ProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting provisioningParameters as " +
				"*aci.ProvisioningParameters",
		)
	}

	outputs, err := s.armDeployer.Deploy(
		pc.ARMDeploymentName,
		standardProvisioningContext.ResourceGroup,
		standardProvisioningContext.Location,
		armTemplateBytes,
		pp, // Go template params
		map[string]interface{}{ // ARM template params
			"name":       pc.ContainerName,
			"image":      pp.ImageName,
			"cpuCores":   pp.NumberCores,
			"memoryInGb": fmt.Sprintf("%f", pp.Memory),
		},
		standardProvisioningContext.Tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	// We don't check if this is ok, because "no public IP" is a legitimate
	// scenario.
	publicIPv4Address, ok := outputs["publicIPv4Address"].(string)
	if ok {
		pc.PublicIPv4Address = publicIPv4Address
	}

	return pc, nil
}
