package rediscache

import (
	"fmt"
)

func (m *module) ValidateBindingParameters(
	bindingParameters interface{},
) error {
	// There are no parameters for binding to Redis, so there is nothing
	// to validate
	return nil
}

func (m *module) Bind(
	provisioningContext interface{},
	bindingParameters interface{},
) (interface{}, interface{}, error) {
	pc, ok := provisioningContext.(*redisProvisioningContext)
	if !ok {
		return nil, nil, fmt.Errorf(
			"error casting provisioningContext as redisProvisioningContext",
		)
	}

	return &redisBindingContext{},
		&redisCredentials{
			Host:     pc.FullyQualifiedDomainName,
			Password: pc.PrimaryKey,
			Port:     6379,
		},
		nil
}
