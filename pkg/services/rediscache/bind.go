package rediscache

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (s *serviceManager) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to Redis, so there is nothing
	// to validate
	return nil
}

func (s *serviceManager) Bind(
	service.Instance,
	service.BindingParameters,
) (service.BindingDetails, error) {
	return &redisBindingDetails{}, nil
}

func (s *serviceManager) GetCredentials(
	instance service.Instance,
	_ service.Binding,
) (service.Credentials, error) {
	dt, ok := instance.Details.(*redisInstanceDetails)
	if !ok {
		return nil, fmt.Errorf(
			"error casting instance.Details as *redisInstanceDetails",
		)
	}
	return &Credentials{
		Host:     dt.FullyQualifiedDomainName,
		Password: dt.PrimaryKey,
		Port:     6379,
	}, nil
}
