package storage

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *serviceManager) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to Storage, so there is nothing
	// to validate
	return nil
}

func (s *serviceManager) Bind(
	service.Instance,
	service.BindingParameters,
) (service.BindingDetails, error) {
	return &storageBindingDetails{}, nil
}

func (s *serviceManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	dt, ok := instance.Details.(*storageInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *storageInstanceDetails",
		)
	}
	return &Credentials{
		StorageAccountName: dt.StorageAccountName,
		AccessKey:          dt.AccessKey,
		ContainerName:      dt.ContainerName,
	}, nil
}
