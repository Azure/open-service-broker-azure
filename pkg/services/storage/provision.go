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
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instanceDetails{
		ARMDeploymentName:  uuid.NewV4().String(),
		StorageAccountName: generate.NewIdentifier(),
	}

	storeKind, ok := instance.Plan.
		GetProperties().Extended[kindKey].(storageKind)
	if !ok {
		return nil, errors.New(
			"error retrieving the storage kind from the plan",
		)
	}

	// Add context that is specific to certain plans
	switch storeKind {
	case storageKindBlobContainer:
		dt.ContainerName = uuid.NewV4().String()
	}
	return &dt, nil
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*instanceDetails)
	storeKind, ok := instance.Plan.GetProperties().Extended[kindKey].(storageKind)
	if !ok {
		return nil, errors.New(
			"error retrieving the storage kind from the plan",
		)
	}

	location := instance.ProvisioningParameters.GetString("location")
	goTemplateParams := map[string]interface{}{
		"name":     dt.StorageAccountName,
		"location": location,
	}
	switch storeKind {
	case storageKindGeneralPurposeStorageAcccount:
		goTemplateParams["kind"] = "Storage"
	case storageKindGeneralPurposeV2StorageAccount:
		goTemplateParams["kind"] = "StorageV2"
		goTemplateParams["accessTier"] = "Hot"
	case storageKindBlobStorageAccount, storageKindBlobContainer:
		goTemplateParams["kind"] = "BlobStorage"
		goTemplateParams["accessTier"] = "Hot"
	}

	tagsObj := instance.ProvisioningParameters.GetObject("tags")
	tags := make(map[string]string, len(tagsObj.Data))
	for k := range tagsObj.Data {
		tags[k] = tagsObj.GetString(k)
	}

	outputs, err := s.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ProvisioningParameters.GetString("resourceGroup"),
		instance.ProvisioningParameters.GetString("location"),
		armTemplateBytes,
		goTemplateParams,         // Go template params
		map[string]interface{}{}, // ARM template params
		tags,
	)
	if err != nil {
		return nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	dt.AccessKey, ok = outputs["accessKey"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"error retrieving primary access key from deployment: %s",
			err,
		)
	}
	return dt, nil
}

func (s *serviceManager) createBlobContainer(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*instanceDetails)

	client, _ := storage.NewBasicClient(dt.StorageAccountName, dt.AccessKey)
	blobCli := client.GetBlobService()
	container := blobCli.GetContainerReference(dt.ContainerName)
	options := storage.CreateContainerOptions{
		Access: storage.ContainerAccessTypePrivate,
	}
	_, err := container.CreateIfNotExists(&options)
	if err != nil {
		return nil, errors.New(
			"error creating container",
		)
	}
	return instance.Details, nil
}
