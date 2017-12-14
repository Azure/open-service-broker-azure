package memory

import (
	"github.com/Azure/open-service-broker-azure/pkg/crypto"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/storage"
)

type store struct {
	codec     crypto.Codec
	instances map[string][]byte
	bindings  map[string][]byte
}

// NewStore returns a new memory-based implementation of the storage.Store used
// for testing
func NewStore(codec crypto.Codec) storage.Store {
	return &store{
		codec:     codec,
		instances: make(map[string][]byte),
		bindings:  make(map[string][]byte),
	}
}

func (s *store) WriteInstance(instance service.Instance) error {
	json, err := instance.ToJSON(s.codec)
	if err != nil {
		return err
	}
	s.instances[instance.InstanceID] = json
	return nil
}

func (s *store) GetInstance(
	instanceID string,
	pp service.ProvisioningParameters,
	up service.UpdatingParameters,
	pc service.ProvisioningContext,
) (
	service.Instance,
	bool,
	error,
) {
	json, ok := s.instances[instanceID]
	if !ok {
		return service.Instance{}, false, nil
	}
	instance, err := service.NewInstanceFromJSON(json, pp, up, pc, s.codec)
	return instance, err == nil, err
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
	json, err := binding.ToJSON(s.codec)
	if err != nil {
		return err
	}
	s.bindings[binding.BindingID] = json
	return nil
}

func (s *store) GetBinding(
	bindingID string,
	bp service.BindingParameters,
	bc service.BindingContext,
	cr service.Credentials,
) (service.Binding, bool, error) {
	json, ok := s.bindings[bindingID]
	if !ok {
		return service.Binding{}, false, nil
	}
	binding, err := service.NewBindingFromJSON(json, bp, bc, cr, s.codec)
	return binding, err == nil, err
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
