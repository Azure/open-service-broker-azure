package storage

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	uuid "github.com/satori/go.uuid"
)

func (b *blobContainerManager) GetProvisioner(
	_ service.Plan,
) (service.Provisioner, error) {
	return service.NewProvisioner(
		service.NewProvisioningStep("preProvision", b.preProvision),
		service.NewProvisioningStep("createBlobContainer", b.createBlobContainer),
	)
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
