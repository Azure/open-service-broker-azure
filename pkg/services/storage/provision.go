package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/Azure/open-service-broker-azure/pkg/generate"
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
	plan service.Plan,
) (service.Provisioner, error) {
	provisioningSteps := []service.ProvisioningStep{
		service.NewProvisioningStep("preProvision", s.preProvision),
		service.NewProvisioningStep("deployARMTemplate", s.deployARMTemplate),
	}

	storeKind, ok := plan.GetProperties().Extended[kindKey].(storageKind)
	if !ok {
		return nil, errors.New(
			"error retrieving the storage kind from the plan: %s",
		)
	}

	// Add provisioning steps that are specific to certain plans
	switch storeKind {
	case storageKindBlobContainer:
		provisioningSteps = append(
			provisioningSteps,
			service.NewProvisioningStep("createBlobContainer", s.createBlobContainer),
		)
	}

	return service.NewProvisioner(provisioningSteps...)
}

func (s *serviceManager) preProvision(
	_ context.Context,
	_ string, // instanceID
	plan service.Plan,
	_ service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	_ service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*storageProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *storageProvisioningContext",
		)
	}
	pc.ARMDeploymentName = uuid.NewV4().String()
	pc.StorageAccountName = generate.NewIdentifier()

	storeKind, ok := plan.GetProperties().Extended[kindKey].(storageKind)
	if !ok {
		return nil, errors.New(
			"error retrieving the storage kind from the plan",
		)
	}

	// Add context that is specific to certain plans
	switch storeKind {
	case storageKindBlobContainer:
		pc.ContainerName = uuid.NewV4().String()
	}

	return pc, nil
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	_ string, // instanceID
	plan service.Plan,
	standardProvisioningContext service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	_ service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*storageProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *storageProvisioningContext",
		)
	}
	storeKind, ok := plan.GetProperties().Extended[kindKey].(storageKind)
	if !ok {
		return nil, errors.New(
			"error retrieving the storage kind from the plan",
		)
	}

	var armTemplateBytes []byte
	switch storeKind {
	case storageKindGeneralPurposeStorageAcccount:
		armTemplateBytes = armTemplateBytesGeneralPurposeStorage
	case storageKindBlobStorageAccount, storageKindBlobContainer:
		armTemplateBytes = armTemplateBytesBlobStorage
	}
	armTemplateParameters := map[string]interface{}{
		"name": pc.StorageAccountName,
	}
	outputs, err := s.armDeployer.Deploy(
		pc.ARMDeploymentName,
		standardProvisioningContext.ResourceGroup,
		standardProvisioningContext.Location,
		armTemplateBytes,
		nil, // Go template params
		armTemplateParameters, // ARM template params
		standardProvisioningContext.Tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	pc.AccessKey, ok = outputs["accessKey"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving primary access key from deployment: %s",
			err,
		)
	}

	return pc, nil
}

func (s *serviceManager) createBlobContainer(
	_ context.Context,
	_ string, // instanceID
	_ service.Plan,
	_ service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	_ service.ProvisioningParameters,
) (service.ProvisioningContext, error) {
	pc, ok := provisioningContext.(*storageProvisioningContext)
	if !ok {
		return nil, errors.New(
			"error casting provisioningContext as *storageProvisioningContext",
		)
	}

	client, _ := storage.NewBasicClient(pc.StorageAccountName, pc.AccessKey)
	blobCli := client.GetBlobService()
	container := blobCli.GetContainerReference(pc.ContainerName)
	options := storage.CreateContainerOptions{
		Access: storage.ContainerAccessTypePrivate,
	}
	_, err := container.CreateIfNotExists(&options)
	if err != nil {
		return nil, errors.New(
			"error creating container",
		)
	}

	return pc, nil
}
