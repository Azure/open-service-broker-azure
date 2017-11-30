package rediscache

import (
	"fmt"

	"github.com/Azure/azure-service-broker/pkg/service"
)

func (s *serviceManager) ValidateBindingParameters(
	bindingParameters service.BindingParameters,
) error {
	// There are no parameters for binding to Redis, so there is nothing
	// to validate
	return nil
}

func (s *serviceManager) Bind(
	_ service.StandardProvisioningContext,
	provisioningContext service.ProvisioningContext,
	bindingParameters service.BindingParameters,
) (service.BindingContext, service.Credentials, error) {
	pc, ok := provisioningContext.(*redisProvisioningContext)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting provisioningContext as *redisProvisioningContext",
		)
	}

	return &redisBindingContext{},
		&Credentials{
			Host:     pc.FullyQualifiedDomainName,
			Password: pc.PrimaryKey,
			Port:     6379,
		},
		nil
}
