package fake

import (
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/storage"
)

type store struct {
}

// NewStore returns a new fake implementation of the storage.Store used for
// testing
func NewStore() storage.Store {
	return &store{}
}

func (s *store) WriteInstance(instance *service.Instance) error {
	return nil
}

func (s *store) GetInstance(instanceID string) (*service.Instance, bool, error) {
	return nil, false, nil
}

func (s *store) DeleteInstance(instanceID string) (bool, error) {
	return false, nil
}

func (s *store) WriteBinding(binding *service.Binding) error {
	return nil
}

func (s *store) GetBinding(bindingID string) (*service.Binding, bool, error) {
	return nil, false, nil
}

func (s *store) DeleteBinding(bindingID string) (bool, error) {
	return false, nil
}

func (s *store) TestConnection() error {
	return nil
}
