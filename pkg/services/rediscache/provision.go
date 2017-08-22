package rediscache

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/azure"
	"github.com/Azure/azure-service-broker/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (m *module) ValidateProvisioningParameters(
	provisioningParameters interface{},
) error {
	pp, ok := provisioningParameters.(*redisProvisioningParameters)
	if !ok {
		return errors.New(
			"error casting provisioningParameters as " +
				"redisProvisioningParameters",
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
	provisioningContext interface{},
	provisioningParameters interface{}, // nolint: unparam
) (interface{}, error) {
	pc, ok := provisioningContext.(*redisProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as redisProvisioningContext",
		)
	}

	// generate names for resources
	pc.ResourceGroupName = uuid.NewV4().String()
	pc.ARMDeploymentName = uuid.NewV4().String()
	pc.ServerName = uuid.NewV4().String()
	return pc, nil
}

func (m *module) deployARMTemplate(
	ctx context.Context, // nolint: unparam
	instanceID string, // nolint: unparam
	serviceID string,
	planID string,
	provisioningContext interface{},
	provisioningParameters interface{},
) (interface{}, error) {
	pc, ok := provisioningContext.(*redisProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as redisProvisioningContext",
		)
	}
	pp, ok := provisioningParameters.(*redisProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting provisioningParameters as " +
				"redisProvisioningParameters",
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
		// TODO: Values in this map should vary according to the serviceID and planID
		// selected
		map[string]interface{}{
			"serverName":         pc.ServerName,
			"redisCacheSKU":      plan.GetProperties().Extended["redisCacheSKU"],
			"redisCacheFamily":   plan.GetProperties().Extended["redisCacheFamily"],
			"redisCacheCapacity": plan.GetProperties().Extended["redisCacheCapacity"],
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	fullyQualifiedDomainName, ok := outputs["fullyQualifiedDomainName"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving fully qualified domain name from deployment: %s",
			err,
		)
	}
	pc.FullyQualifiedDomainName = fullyQualifiedDomainName

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
