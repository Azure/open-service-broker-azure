package storage

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *serviceManager) ValidateBindingParameters(
	service.BindingParameters,
	service.SecureBindingParameters,
) error {
	// There are no parameters for binding to Storage, so there is nothing
	// to validate
	return nil
}

func (s *serviceManager) Bind(
	service.Instance,
	service.BindingParameters,
	service.SecureBindingParameters,
) (service.BindingDetails, service.SecureBindingDetails, error) {
	return &storageBindingDetails{}, &storageSecureBindingDetails{}, nil
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
	sdt, ok := instance.SecureDetails.(*storageSecureInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.SecureDetails as *storageSecureInstanceDetails",
		)
	}
	return &Credentials{
		StorageAccountName: dt.StorageAccountName,
		AccessKey:          sdt.AccessKey,
		ContainerName:      dt.ContainerName,
	}, nil
}
