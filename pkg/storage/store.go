package storage

import (
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/go-redis/redis"
)

// Store is an interface to be implemented by types capable of handling
// persistence for other broker-related types
type Store interface {
	// WriteInstance persists the given instance to the underlying storage
	WriteInstance(instance *service.Instance) error
	// GetInstance retrieves a persisted instance from the underlying storage by
	// instance id
	GetInstance(instanceID string) (*service.Instance, bool, error)
	// DeleteInstance deletes a persisted instance from the underlying storage by
	// instance id
	DeleteInstance(instanceID string) (bool, error)
	// WriteBinding persists the given binding to the underlying storage
	WriteBinding(binding *service.Binding) error
	// GetBinding retrieves a persisted instance from the underlying storage by
	// binding id
	GetBinding(bindingID string) (*service.Binding, bool, error)
	// DeleteBinding deletes a persisted binding from the underlying storage by
	// binding id
	DeleteBinding(bindingID string) (bool, error)
}

type store struct {
	redisClient *redis.Client
}

// NewStore returns a new Redis-based implementation of the Store interface
func NewStore(redisClient *redis.Client) Store {
	return &store{
		redisClient: redisClient,
	}
}

func (s *store) WriteInstance(instance *service.Instance) error {
	json, err := instance.ToJSONString()
	if err != nil {
		return err
	}
	return s.redisClient.Set(instance.InstanceID, json, 0).Err()
}

func (s *store) GetInstance(instanceID string) (*service.Instance, bool, error) {
	strCmd := s.redisClient.Get(instanceID)
	err := strCmd.Err()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	bytes, err := strCmd.Bytes()
	if err != nil {
		return nil, false, err
	}
	instance, err := service.NewInstanceFromJSONString(string(bytes))
	if err != nil {
		return nil, false, err
	}
	return instance, true, nil
}

func (s *store) DeleteInstance(instanceID string) (bool, error) {
	strCmd := s.redisClient.Get(instanceID)
	err := strCmd.Err()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	err = s.redisClient.Del(instanceID).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *store) WriteBinding(binding *service.Binding) error {
	json, err := binding.ToJSONString()
	if err != nil {
		return err
	}
	return s.redisClient.Set(binding.BindingID, json, 0).Err()
}

func (s *store) GetBinding(bindingID string) (*service.Binding, bool, error) {
	strCmd := s.redisClient.Get(bindingID)
	err := strCmd.Err()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	bytes, err := strCmd.Bytes()
	if err != nil {
		return nil, false, err
	}
	binding, err := service.NewBindingFromJSONString(string(bytes))
	if err != nil {
		return nil, false, err
	}
	return binding, true, nil
}

func (s *store) DeleteBinding(bindingID string) (bool, error) {
	strCmd := s.redisClient.Get(bindingID)
	err := strCmd.Err()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	err = s.redisClient.Del(bindingID).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}
