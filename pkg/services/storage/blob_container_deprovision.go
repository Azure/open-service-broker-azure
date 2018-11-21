package storage

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (b *blobContainerManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteBlobContainer", b.deleteBlobContainer),
	)
}

func (b *blobContainerManager) deleteBlobContainer(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, error) {
	dt := instance.Details.(*instanceDetails)
	if err := deleteBlobContainer(
		dt.StorageAccountName,
		dt.AccessKey,
		dt.ContainerName,
	); err != nil {
		return nil, err
	}
	return instance.Details, nil
}
