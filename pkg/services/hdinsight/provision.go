package hdinsight

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/azure"
	"github.com/Azure/azure-service-broker/pkg/generate"
	"github.com/Azure/azure-service-broker/pkg/service"
	az "github.com/Azure/go-autorest/autorest/azure"
	uuid "github.com/satori/go.uuid"
)

func (m *module) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	pp, ok := provisioningParameters.(*ProvisioningParameters)
	if !ok {
		return errors.New(
			"error casting provisioningParameters as " +
				"hdinsightProvisioningParameters",
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
	provisioningParameters service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*hdinsightProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as hdinsightProvisioningContext",
		)
	}
	pp, ok := provisioningParameters.(*ProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting provisioningParameters as " +
				"*hdinsight.ProvisioningParameters",
		)
	}

	if pp.ResourceGroup != "" {
		pc.ResourceGroupName = pp.ResourceGroup
	} else {
		pc.ResourceGroupName = uuid.NewV4().String()
	}
	pc.ARMDeploymentName = uuid.NewV4().String()
	pc.ClusterName = "hd-" + uuid.NewV4().String()
	pc.ClusterLoginUserName = generate.NewIdentifier()
	pc.ClusterLoginPassword = generate.NewPassword()
	pc.SSHUserName = generate.NewIdentifier()
	pc.SSHPassword = generate.NewPassword()
	pc.StorageAccountName = generate.NewIdentifier()
	pc.StorageAccountKey = generate.NewPassword()

	// Get BlobStorageEndpoint
	azureConfig, err := azure.GetConfig()
	if err != nil {
		return nil, err
	}
	azureEnvironment, err := az.EnvironmentFromName(azureConfig.Environment)
	if err != nil {
		return nil, err
	}
	storageEndpointSuffix := azureEnvironment.StorageEndpointSuffix
	pc.BlobStorageEndpoint = fmt.Sprintf(
		"%s.blob.%s",
		pc.StorageAccountName,
		storageEndpointSuffix,
	)

	pc.BlobStorageContainerName = generate.NewIdentifier()
	return pc, nil
}

func (m *module) deployARMTemplate(
	ctx context.Context, // nolint: unparam
	instanceID string, // nolint: unparam
	serviceID string, // nolint: unparam
	planID string, // nolint: unparam
	provisioningContext service.ProvisioningContext,
	provisioningParameters service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*hdinsightProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as hdinsightProvisioningContext",
		)
	}
	pp, ok := provisioningParameters.(*ProvisioningParameters)
	if !ok {
		return nil, errors.New(
			"error casting provisioningParameters as " +
				"hdinsightProvisioningParameters",
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
	planName := plan.GetProperties().Name

	armParams := map[string]interface{}{
		"clusterName":              pc.ClusterName,
		"clusterLoginUserName":     pc.ClusterLoginUserName,
		"clusterLoginPassword":     pc.ClusterLoginPassword,
		"sshUserName":              pc.SSHUserName,
		"sshPassword":              pc.SSHPassword,
		"storageAccountName":       pc.StorageAccountName,
		"blobStorageContainerName": pc.BlobStorageContainerName,
		"blobStorageEndpoint":      pc.BlobStorageEndpoint,
	}
	if pp.ClusterWorkerNodeCount != 0 {
		armParams["clusterWorkerNodeCount"] = pp.ClusterWorkerNodeCount
	}
	outputs, err := m.armDeployer.Deploy(
		pc.ARMDeploymentName,
		pc.ResourceGroupName,
		pp.Location,
		armTemplateBytes[planName],
		armParams,
		pp.Tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	storageAccountKey, ok := outputs["storageAccountKey"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving storage account key from deployment: %s",
			err,
		)
	}
	pc.StorageAccountKey = storageAccountKey

	return pc, nil
}
