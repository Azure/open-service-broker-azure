// +build experimental

package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (c *cosmosAccountManager) Unbind(
	_ service.Instance,
	_ service.Binding,
) error {
	return nil
}
