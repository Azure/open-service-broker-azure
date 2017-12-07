package hdinsight

import (
	"context"
	"errors"
	"fmt"

	az "github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/open-service-broker-azure/pkg/azure"
	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (s *serviceManager) ValidateProvisioningParameters(
	provisioningParameters service.ProvisioningParameters,
) error {
	if _, ok := provisioningParameters.(*ProvisioningParameters); !ok {
		return errors.New(
			"error casting provisioningParameters as " +
				"hdinsightProvisioningParameters",
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
	pc, ok := provisioningContext.(*hdinsightProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as hdinsightProvisioningContext",
		)
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

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	_ string, // instanceID
	plan service.Plan,
	standardProvisioningContext service.StandardProvisioningContext,
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
	outputs, err := s.armDeployer.Deploy(
		pc.ARMDeploymentName,
		standardProvisioningContext.ResourceGroup,
		standardProvisioningContext.Location,
		armTemplateBytes[planName],
		nil, // Go template params
		armParams,
		standardProvisioningContext.Tags,
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
