package storage

import (
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/crypto"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/go-redis/redis"
)

// Store is an interface to be implemented by types capable of handling
// persistence for other broker-related types
type Store interface {
	// WriteInstance persists the given instance to the underlying storage
	WriteInstance(instance service.Instance) error
	// GetInstance retrieves a persisted instance from the underlying storage by
	// instance id
	GetInstance(instanceID string) (service.Instance, bool, error)
	// DeleteInstance deletes a persisted instance from the underlying storage by
	// instance id
	DeleteInstance(instanceID string) (bool, error)
	// WriteBinding persists the given binding to the underlying storage
	WriteBinding(binding service.Binding) error
	// GetBinding retrieves a persisted instance from the underlying storage by
	// binding id
	GetBinding(bindingID string) (service.Binding, bool, error)
	// DeleteBinding deletes a persisted binding from the underlying storage by
	// binding id
	DeleteBinding(bindingID string) (bool, error)
	// TestConnection tests the connection to the underlying database (if there
	// is one)
	TestConnection() error
}

type store struct {
	redisClient *redis.Client
	catalog     service.Catalog
	codec       crypto.Codec
}

// NewStore returns a new Redis-based implementation of the Store interface
func NewStore(
	redisClient *redis.Client,
	catalog service.Catalog,
	codec crypto.Codec,
) Store {
	return &store{
		redisClient: redisClient,
		catalog:     catalog,
		codec:       codec,
	}
}

func (s *store) WriteInstance(instance service.Instance) error {
	json, err := instance.ToJSON(s.codec)
	if err != nil {
		return err
	}
	return s.redisClient.Set(instance.InstanceID, json, 0).Err()
}

func (s *store) GetInstance(instanceID string) (service.Instance, bool, error) {
	strCmd := s.redisClient.Get(instanceID)
	if err := strCmd.Err(); err == redis.Nil {
		return service.Instance{}, false, nil
	} else if err != nil {
		return service.Instance{}, false, err
	}
	bytes, err := strCmd.Bytes()
	if err != nil {
		return service.Instance{}, false, err
	}
	instance, err := service.NewInstanceFromJSON(bytes, nil, nil, nil, s.codec)
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
	serviceManager := svc.GetServiceManager()
	instance, err = service.NewInstanceFromJSON(
		bytes,
		serviceManager.GetEmptyProvisioningParameters(),
		serviceManager.GetEmptyUpdatingParameters(),
		serviceManager.GetEmptyProvisioningContext(),
		s.codec,
	)
	return instance, err == nil, err
}

func (s *store) DeleteInstance(instanceID string) (bool, error) {
	strCmd := s.redisClient.Get(instanceID)
	if err := strCmd.Err(); err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	if err := s.redisClient.Del(instanceID).Err(); err != nil {
		return false, err
	}
	return true, nil
}

func (s *store) WriteBinding(binding service.Binding) error {
	json, err := binding.ToJSON(s.codec)
	if err != nil {
		return err
	}
	return s.redisClient.Set(binding.BindingID, json, 0).Err()
}

func (s *store) GetBinding(bindingID string) (service.Binding, bool, error) {
	strCmd := s.redisClient.Get(bindingID)
	if err := strCmd.Err(); err == redis.Nil {
		return service.Binding{}, false, nil
	} else if err != nil {
		return service.Binding{}, false, err
	}
	bytes, err := strCmd.Bytes()
	if err != nil {
		return service.Binding{}, false, err
	}
	binding, err := service.NewBindingFromJSON(bytes, nil, nil, nil, s.codec)
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
		bytes,
		serviceManager.GetEmptyBindingParameters(),
		serviceManager.GetEmptyBindingContext(),
		serviceManager.GetEmptyCredentials(),
		s.codec,
	)
	return binding, err == nil, err
}

func (s *store) DeleteBinding(bindingID string) (bool, error) {
	strCmd := s.redisClient.Get(bindingID)
	if err := strCmd.Err(); err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	if err := s.redisClient.Del(bindingID).Err(); err != nil {
		return false, err
	}
	return true, nil
}

func (s *store) TestConnection() error {
	return s.redisClient.Ping().Err()
}
