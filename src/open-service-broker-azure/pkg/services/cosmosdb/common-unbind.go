// +build experimental

package cosmosdb

import (
	"open-service-broker-azure/pkg/service"
)

func (c *cosmosAccountManager) Unbind(
	_ service.Instance,
	_ service.Binding,
) error {
	return nil
}
