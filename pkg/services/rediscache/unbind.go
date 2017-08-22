package rediscache

import "fmt"

func (m *module) Unbind(
	provisioningContext interface{}, // nolint: unparam
	bindingContext interface{},
) error {
	pc, ok := provisioningContext.(*redisProvisioningContext)
	if !ok {
		return fmt.Errorf(
			"error casting provisioningContext as redisProvisioningContext",
		)
	}

	err := getDBConnection(pc)

	return err
}
