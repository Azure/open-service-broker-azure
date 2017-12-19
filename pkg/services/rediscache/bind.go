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
	instance service.Instance,
	_ service.BindingParameters,
) (service.BindingContext, service.Credentials, error) {
	dt, ok := instance.Details.(*redisInstanceDetails)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting instance.Details as *redisInstanceDetails",
		)
	}

	return &redisBindingContext{},
		&Credentials{
			Host:     dt.FullyQualifiedDomainName,
			Password: dt.PrimaryKey,
			Port:     6379,
		},
		nil
}
