package cosmosdb

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (m *mongoAccountManager) Unbind(
	_ service.Instance,
	_ service.Binding,
) error {
	return nil
}
