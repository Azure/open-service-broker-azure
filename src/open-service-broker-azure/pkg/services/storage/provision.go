// +build experimental

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
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := instanceDetails{
		ARMDeploymentName:  uuid.NewV4().String(),
		StorageAccountName: generate.NewIdentifier(),
	}

	storeKind, ok := instance.Plan.
		GetProperties().Extended[kindKey].(storageKind)
	if !ok {
		return nil, nil, errors.New(
			"error retrieving the storage kind from the plan",
		)
	}

	// Add context that is specific to certain plans
	switch storeKind {
	case storageKindBlobContainer:
		dt.ContainerName = uuid.NewV4().String()
	}

	dtMap, err := service.GetMapFromStruct(dt)
	return dtMap, instance.SecureDetails, err
}

func (s *serviceManager) deployARMTemplate(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := instanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}
	sdt := secureInstanceDetails{}
	if err := service.GetStructFromMap(instance.SecureDetails, &sdt); err != nil {
		return nil, nil, err
	}
	storeKind, ok := instance.Plan.GetProperties().Extended[kindKey].(storageKind)
	if !ok {
		return nil, nil, errors.New(
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
		"name": dt.StorageAccountName,
	}
	outputs, err := s.armDeployer.Deploy(
		dt.ARMDeploymentName,
		instance.ResourceGroup,
		instance.Location,
		armTemplateBytes,
		nil, // Go template params
		armTemplateParameters, // ARM template params
		instance.Tags,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("error deploying ARM template: %s", err)
	}

	sdt.AccessKey, ok = outputs["accessKey"].(string)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error retrieving primary access key from deployment: %s",
			err,
		)
	}

	sdtMap, err := service.GetMapFromStruct(sdt)
	return instance.Details, sdtMap, err
}

func (s *serviceManager) createBlobContainer(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	dt := instanceDetails{}
	if err := service.GetStructFromMap(instance.Details, &dt); err != nil {
		return nil, nil, err
	}
	sdt := secureInstanceDetails{}
	if err := service.GetStructFromMap(instance.SecureDetails, &sdt); err != nil {
		return nil, nil, err
	}

	client, _ := storage.NewBasicClient(dt.StorageAccountName, sdt.AccessKey)
	blobCli := client.GetBlobService()
	container := blobCli.GetContainerReference(dt.ContainerName)
	options := storage.CreateContainerOptions{
		Access: storage.ContainerAccessTypePrivate,
	}
	_, err := container.CreateIfNotExists(&options)
	if err != nil {
		return nil, nil, errors.New(
			"error creating container",
		)
	}

	return instance.Details, instance.SecureDetails, nil
}
