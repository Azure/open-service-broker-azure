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
	instance service.Instance,
	_ service.BindingParameters,
) (service.BindingContext, service.Credentials, error) {
	dt, ok := instance.Details.(*storageInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *storageInstanceDetails",
		)
	}

	return &storageBindingContext{},
		&Credentials{
			StorageAccountName: dt.StorageAccountName,
			AccessKey:          dt.AccessKey,
			ContainerName:      dt.ContainerName,
		},
		nil
}
