package memory

import (
	"fmt"
	"sync"

	"github.com/Azure/open-service-broker-azure/pkg/crypto"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/storage"
)

type store struct {
	catalog                       service.Catalog
	codec                         crypto.Codec
	instances                     map[string][]byte
	instanceAliases               map[string]string
	bindings                      map[string][]byte
	instanceAliasChildCounts      map[string]int64
	instanceAliasChildCountsMutex sync.Mutex
}

// NewStore returns a new memory-based implementation of the storage.Store used
// for testing
func NewStore(catalog service.Catalog, codec crypto.Codec) storage.Store {
	return &store{
		catalog:                  catalog,
		codec:                    codec,
		instances:                make(map[string][]byte),
		instanceAliases:          make(map[string]string),
		bindings:                 make(map[string][]byte),
		instanceAliasChildCounts: make(map[string]int64),
	}
}

func (s *store) WriteInstance(instance service.Instance) error {
	json, err := instance.ToJSON(s.codec)
	if err != nil {
		return err
	}
	s.instances[instance.InstanceID] = json
	if instance.Alias != "" {
		s.instanceAliases[instance.Alias] = instance.InstanceID
	}
	if instance.ParentAlias != "" {
		s.instanceAliasChildCountsMutex.Lock()
		defer s.instanceAliasChildCountsMutex.Unlock()
		s.instanceAliasChildCounts[instance.ParentAlias]++
	}
	return nil
}

func (s *store) GetInstance(instanceID string) (
	service.Instance,
	bool,
	error,
) {
	json, ok := s.instances[instanceID]
	if !ok {
		return service.Instance{}, false, nil
	}
	instance, err := service.NewInstanceFromJSON(
		json,
		nil,
		nil,
		nil,
		nil,
		nil,
		s.codec,
	)
	if err != nil {
		return instance, false, err
	}
	svc, ok := s.catalog.GetService(instance.ServiceID)
	if !ok {
		return instance,
			false,
			fmt.Errorf(
				`service not found in catalog for service ID "%s"`,
				instance.ServiceID,
			)
	}
	plan, ok := svc.GetPlan(instance.PlanID)
	if !ok {
		return instance,
			false,
			fmt.Errorf(
				`plan not found for planID "%s" for service "%s" in the catalog`,
				instance.PlanID,
				instance.ServiceID,
			)
	}
	serviceManager := svc.GetServiceManager()
	instance, err = service.NewInstanceFromJSON(
		json,
		serviceManager.GetEmptyProvisioningParameters(),
		serviceManager.GetEmptySecureProvisioningParameters(),
		serviceManager.GetEmptyUpdatingParameters(),
		serviceManager.GetEmptyInstanceDetails(),
		serviceManager.GetEmptySecureInstanceDetails(),
		s.codec,
	)
	instance.Service = svc
	instance.Plan = plan
	return instance, err == nil, err
}

func (s *store) GetInstanceByAlias(alias string) (
	service.Instance,
	bool,
	error,
) {
	instanceID, ok := s.instanceAliases[alias]
	if !ok {
		return service.Instance{}, false, nil
	}
	return s.GetInstance(instanceID)
}

func (s *store) DeleteInstance(instanceID string) (bool, error) {
	instance, ok, err := s.GetInstance(instanceID)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	delete(s.instances, instanceID)
	if instance.Alias != "" {
		delete(s.instanceAliases, instance.Alias)
	}
	if instance.ParentAlias != "" {
		s.instanceAliasChildCountsMutex.Lock()
		defer s.instanceAliasChildCountsMutex.Unlock()
		s.instanceAliasChildCounts[instance.ParentAlias]--
	}
	return true, nil
}

func (s *store) GetInstanceChildCountByAlias(alias string) (int64, error) {
	s.instanceAliasChildCountsMutex.Lock()
	defer s.instanceAliasChildCountsMutex.Unlock()
	return s.instanceAliasChildCounts[alias], nil
}

func (s *store) WriteBinding(binding service.Binding) error {
	json, err := binding.ToJSON(s.codec)
	if err != nil {
		return err
	}
	s.bindings[binding.BindingID] = json
	return nil
}

func (s *store) GetBinding(bindingID string) (service.Binding, bool, error) {
	json, ok := s.bindings[bindingID]
	if !ok {
		return service.Binding{}, false, nil
	}
	binding, err := service.NewBindingFromJSON(json, nil, nil, nil, s.codec)
	if err != nil {
		return binding, false, err
	}
	svc, ok := s.catalog.GetService(binding.ServiceID)
	if !ok {
		return binding,
			false,
			fmt.Errorf(
				`service not found in catalog for service ID "%s"`,
				binding.ServiceID,
			)
	}
	serviceManager := svc.GetServiceManager()
	binding, err = service.NewBindingFromJSON(
		json,
		serviceManager.GetEmptyBindingParameters(),
		serviceManager.GetEmptyBindingDetails(),
		serviceManager.GetEmptySecureBindingDetails(),
		s.codec,
	)
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
