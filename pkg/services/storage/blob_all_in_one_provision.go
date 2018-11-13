package storage

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/storage"
	"github.com/Azure/open-service-broker-azure/pkg/generate"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (b *blobAllInOneManager) GetProvisioner(
	plan service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", b.preProvision),
		service.NewProvisioningStep("deployARMTemplate", b.deployARMTemplate),
		service.NewProvisioningStep("createBlobContainer", b.createBlobContainer),
	)
}

func (b *blobAllInOneManager) preProvision(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instanceDetails{
		ARMDeploymentName:  uuid.NewV4().String(),
		StorageAccountName: generate.NewIdentifier(),
		ContainerName:      uuid.NewV4().String(),
	}
	return &dt, nil
}

func (b *blobAllInOneManager) createBlobContainer(
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
