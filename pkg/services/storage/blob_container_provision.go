package storage

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (b *blobContainerManager) GetProvisioner(
	_ service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("checkNameAvailability", b.checkNameAvailability),
		service.NewProvisioningStep("preProvision", b.preProvision),
		service.NewProvisioningStep("createBlobContainer", b.createBlobContainer),
	)
}

func (b *blobContainerManager) checkNameAvailability(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	containerName := instance.ProvisioningParameters.GetString("containerName")
	if containerName == "" {
		return nil, nil
	}

	pdt := instance.Parent.Details.(instanceDetails)
	client, _ := storage.NewBasicClient(
		pdt.StorageAccountName,
		pdt.AccessKey,
	)
	blobCli := client.GetBlobService()
	response, err := blobCli.ListContainers(storage.ListContainersParameters{})
	if err != nil {
		return nil, fmt.Errorf("error checking name availability %s", err)
	}
	containers := response.Containers
	for _, container := range containers {
		if containerName == container.Name {
			return nil, fmt.Errorf(
				"container having name %s already exists in the storage account",
				containerName,
			)
		}
	}
	return nil, nil
}

func (b *blobContainerManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	pdt := instance.Parent.Details.(*instanceDetails)
	dt := *pdt
	if instance.ProvisioningParameters.GetString("containerName") != "" {
		dt.ContainerName = instance.ProvisioningParameters.GetString("containerName")
	} else {
		dt.ContainerName = uuid.NewV4().String()
	}
	return &dt, nil
}

func (b *blobContainerManager) createBlobContainer(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*instanceDetails)
	if err := createBlobContainer(
		dt.StorageAccountName,
		dt.AccessKey,
		dt.ContainerName,
	); err != nil {
		return nil, err
	}
	return instance.Details, nil
}
