package storage

import (
	"crypto/tls"
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
	// GetInstanceByID retrieves a persisted instance from the underlying storage
	// by alias
	GetInstanceByAlias(alias string) (service.Instance, bool, error)
	// GetInstanceChildCountByAlias returns the number of child instances
	GetInstanceChildCountByAlias(alias string) (int64, error)
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
	catalog service.Catalog,
	codec crypto.Codec,
	config Config,
) Store {
	redisOpts := &redis.Options{
		Addr:       fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
		Password:   config.RedisPassword,
		DB:         config.RedisDB,
		MaxRetries: 5,
	}
	if config.RedisEnableTLS {
		redisOpts.TLSConfig = &tls.Config{
			ServerName: config.RedisHost,
		}
	}

	return &store{
		redisClient: redis.NewClient(redisOpts),
		catalog:     catalog,
		codec:       codec,
	}
}

func (s *store) WriteInstance(instance service.Instance) error {
	key := getInstanceKey(instance.InstanceID)
	json, err := instance.ToJSON(s.codec)
	if err != nil {
		return err
	}
	pipeline := s.redisClient.TxPipeline()
	pipeline.Set(key, json, 0)
	if instance.Alias != "" {
		aliasKey := getInstanceAliasKey(instance.Alias)
		pipeline.Set(aliasKey, instance.InstanceID, 0)
	}
	if instance.ParentAlias != "" {
		parentAliasChildrenKey := getInstanceAliasChildrenKey(instance.ParentAlias)
		pipeline.SAdd(parentAliasChildrenKey, instance.InstanceID)
	}
	_, err = pipeline.Exec()
	if err != nil {
		return fmt.Errorf(
			`error writing instance "%s": %s`,
			instance.InstanceID,
			err,
		)
	}
	return err
}

func (s *store) GetInstance(instanceID string) (service.Instance, bool, error) {
	key := getInstanceKey(instanceID)
	strCmd := s.redisClient.Get(key)
	if err := strCmd.Err(); err == redis.Nil {
		return service.Instance{}, false, nil
	} else if err != nil {
		return service.Instance{}, false, err
	}
	bytes, err := strCmd.Bytes()
	if err != nil {
		return service.Instance{}, false, err
	}
	instance, err := service.NewInstanceFromJSON(bytes, s.codec)
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
	instance, err = service.NewInstanceFromJSON(
		bytes,
		s.codec,
	)
	instance.Service = svc
	instance.Plan = plan
	if instance.ParentAlias != "" {
		parent, ok, err := s.GetInstanceByAlias(instance.ParentAlias)
		if err != nil {
			return instance, false, fmt.Errorf(
				`error retrieving parent with alias "%s" for instance "%s"`,
				instance.ParentAlias,
				instance.InstanceID,
			)
		}
		if ok {
			instance.Parent = &parent
		}
	}
	return instance, err == nil, err
}

func (s *store) GetInstanceByAlias(
	alias string,
) (service.Instance, bool, error) {
	key := getInstanceAliasKey(alias)
	strCmd := s.redisClient.Get(key)
	if err := strCmd.Err(); err == redis.Nil {
		return service.Instance{}, false, nil
	} else if err != nil {
		return service.Instance{}, false, err
	}
	instanceID, err := strCmd.Result()
	if err != nil {
		return service.Instance{}, false, err
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
	key := getInstanceKey(instanceID)
	pipeline := s.redisClient.TxPipeline()
	pipeline.Del(key)
	if instance.Alias != "" {
		aliasKey := getInstanceAliasKey(instance.Alias)
		pipeline.Del(aliasKey)
	}
	if instance.ParentAlias != "" {
		parentAliasChildrenKey := getInstanceAliasChildrenKey(instance.ParentAlias)
		pipeline.SRem(parentAliasChildrenKey, instance.InstanceID)
	}
	_, err = pipeline.Exec()
	if err != nil {
		return false, fmt.Errorf(
			`error deleting instance "%s": %s`,
			instance.InstanceID,
			err,
		)
	}
	return true, nil
}

func (s *store) GetInstanceChildCountByAlias(alias string) (int64, error) {
	aliasChildrenKey := getInstanceAliasChildrenKey(alias)
	return s.redisClient.SCard(aliasChildrenKey).Result()
}

func getInstanceKey(instanceID string) string {
	return fmt.Sprintf("instances:%s", instanceID)
}

func getInstanceAliasKey(alias string) string {
	return fmt.Sprintf("instances:aliases:%s", alias)
}

func getInstanceAliasChildrenKey(alias string) string {
	return fmt.Sprintf("instances:aliases:%s:children", alias)
}

func (s *store) WriteBinding(binding service.Binding) error {
	key := getBindingKey(binding.BindingID)
	json, err := binding.ToJSON(s.codec)
	if err != nil {
		return err
	}
	return s.redisClient.Set(key, json, 0).Err()
}

func (s *store) GetBinding(bindingID string) (service.Binding, bool, error) {
	key := getBindingKey(bindingID)
	strCmd := s.redisClient.Get(key)
	if err := strCmd.Err(); err == redis.Nil {
		return service.Binding{}, false, nil
	} else if err != nil {
		return service.Binding{}, false, err
	}
	bytes, err := strCmd.Bytes()
	if err != nil {
		return service.Binding{}, false, err
	}
	binding, err := service.NewBindingFromJSON(bytes, s.codec)
	if err != nil {
		return binding, false, err
	}
	binding, err = service.NewBindingFromJSON(bytes, s.codec)
	return binding, err == nil, err
}

func (s *store) DeleteBinding(bindingID string) (bool, error) {
	key := getBindingKey(bindingID)
	strCmd := s.redisClient.Get(key)
	if err := strCmd.Err(); err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	if err := s.redisClient.Del(key).Err(); err != nil {
		return false, err
	}
	return true, nil
}

func getBindingKey(bindingID string) string {
	return fmt.Sprintf("bindings:%s", bindingID)
}

func (s *store) TestConnection() error {
	return s.redisClient.Ping().Err()
}
