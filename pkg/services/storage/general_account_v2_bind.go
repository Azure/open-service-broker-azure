package storage

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (*generalPurposeV2Manager) Bind(
	service.Instance,
	service.BindingParameters,
) (service.BindingDetails, error) {
	return nil, nil
}

// nolint: lll
func (gpv2m *generalPurposeV2Manager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	dt := instance.Details.(*instanceDetails)
	credential := credentials{
		StorageAccountName:          dt.StorageAccountName,
		AccessKey:                   dt.AccessKey,
		ContainerName:               dt.ContainerName,
		PrimaryBlobServiceEndPoint:  fmt.Sprintf("https://%s.blob.core.windows.net/", dt.StorageAccountName),
		PrimaryTableServiceEndPoint: fmt.Sprintf("https://%s.table.core.windows.net/", dt.StorageAccountName),
		PrimaryFileServiceEndPoint:  fmt.Sprintf("https://%s.file.core.windows.net/", dt.StorageAccountName),
		PrimaryQueueServiceEndPoint: fmt.Sprintf("https://%s.queue.core.windows.net/", dt.StorageAccountName),
	}
	return credential, nil
}
