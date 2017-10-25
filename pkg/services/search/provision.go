package search

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/azure"
	"github.com/Azure/azure-service-broker/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (m *module) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	pp, ok := provisioningParameters.(*ProvisioningParameters)
	if !ok {
		return service.NewValidationError(
			"parameters",
			"error casting provisioningParameters as "+
				"*search.ProvisioningParameters",
		)
	}
	if !azure.IsValidLocation(pp.Location) {
		return service.NewValidationError(
			"location",
			fmt.Sprintf(`invalid location: "%s"`, pp.Location),
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
	ctx context.Context, // nolint: unparam
	instanceID string, // nolint: unparam
	serviceID string, // nolint: unparam
	planID string, // nolint: unparam
	provisioningContext service.ProvisioningContext,
	provisioningParameters service.ProvisioningParameters, // nolint: unparam
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*searchProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *searchProvisioningContext",
		)
	}
	pp, ok := provisioningParameters.(*ProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting provisioningParameters as " +
				"*search.ProvisioningParameters",
		)
	}
	if pp.ResourceGroup != "" {
		pc.ResourceGroupName = pp.ResourceGroup
	} else {
		pc.ResourceGroupName = uuid.NewV4().String()
	}
	pc.ARMDeploymentName = uuid.NewV4().String()
	pc.ServiceName = uuid.NewV4().String()
	return pc, nil
}

func (m *module) deployARMTemplate(
	ctx context.Context, // nolint: unparam
	instanceID string, // nolint: unparam
	serviceID string,
	planID string,
	provisioningContext service.ProvisioningContext,
	provisioningParameters service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*searchProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *searchProvisioningContext",
		)
	}
	pp, ok := provisioningParameters.(*ProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting provisioningParameters as " +
				"*search.ProvisioningParameters",
		)
	}
	catalog, err := m.GetCatalog()
	if err != nil {
		return nil, fmt.Errorf("error retrieving catalog: %s", err)
	}
	service, ok := catalog.GetService(serviceID)
	if !ok {
		return nil, fmt.Errorf(
			`service "%s" not found in the "%s" module catalog`,
			serviceID,
			m.GetName(),
		)
	}
	plan, ok := service.GetPlan(planID)
	if !ok {
		return nil, fmt.Errorf(
			`plan "%s" not found for service "%s"`,
			planID,
			serviceID,
		)
	}

	outputs, err := m.armDeployer.Deploy(
		pc.ARMDeploymentName,
		pc.ResourceGroupName,
		pp.Location,
		armTemplateBytes,
		map[string]interface{}{
			"searchServiceName": pc.ServiceName,
			"searchServiceSku":  plan.GetProperties().Extended["searchServiceSku"],
		},
		pp.Tags,
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
