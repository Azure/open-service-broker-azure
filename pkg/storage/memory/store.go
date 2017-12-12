package memory

import (
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/storage"
)

type store struct {
	instances map[string]service.Instance
	bindings  map[string]service.Binding
}

// NewStore returns a new memory-based implementation of the storage.Store used
// for testing
func NewStore() storage.Store {
	return &store{
		instances: make(map[string]service.Instance),
		bindings:  make(map[string]service.Binding),
	}
}

func (s *store) WriteInstance(instance service.Instance) error {
	s.instances[instance.InstanceID] = instance
	return nil
}

func (s *store) GetInstance(instanceID string) (
	service.Instance,
	bool,
	error,
) {
	instance, ok := s.instances[instanceID]
	return instance, ok, nil
}

func (s *store) DeleteInstance(instanceID string) (bool, error) {
	_, ok := s.instances[instanceID]
	if !ok {
		return false, nil
	}
	delete(s.instances, instanceID)
	return true, nil
}

func (s *store) WriteBinding(binding service.Binding) error {
	s.bindings[binding.BindingID] = binding
	return nil
}

func (s *store) GetBinding(bindingID string) (service.Binding, bool, error) {
	binding, ok := s.bindings[bindingID]
	return binding, ok, nil
}

func (s *store) DeleteBinding(bindingID string) (bool, error) {
	_, ok := s.bindings[bindingID]
	if !ok {
		return false, nil
	}
	delete(s.bindings, bindingID)
	return true, nil
}

func (s *store) TestConnection() error {
	return nil
}
