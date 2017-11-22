package aci

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (m *module) ValidateProvisioningParameters(
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

func (m *module) GetProvisioner(string, string) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", m.preProvision),
		service.NewProvisioningStep("deployARMTemplate", m.deployARMTemplate),
	)
}

func (m *module) preProvision(
	_ context.Context,
	_ string, // instanceID
	_ string, // serviceID
	_ string, // planID
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

func (m *module) deployARMTemplate(
	_ context.Context,
	_ string, // instanceID
	_ string, // serviceID
	_ string, // planID
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

	outputs, err := m.armDeployer.Deploy(
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
