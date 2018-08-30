package redis

import (
	"crypto/tls"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/storage"
	"github.com/go-redis/redis"
)

const (
	instances = "instanceList"
	bindings  = "bindingList"
)

type store struct {
	redisClient *redis.Client
	catalog     service.Catalog
}

// NewStore returns a new Redis-based implementation of the Store interface
func NewStore(
	catalog service.Catalog,
	config Config,
) (storage.Store, error) {
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
	}, nil
}

func (s *store) WriteInstance(instance service.Instance) error {
	key := getInstanceKey(instance.InstanceID)
	json, err := instance.ToJSON()
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
	pipeline.SAdd(instances, key)
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
	instance, err := service.NewInstanceFromJSON(bytes, nil, nil)
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
	pps := plan.GetSchemas().ServiceInstances.ProvisioningParametersSchema
	instance, err = service.NewInstanceFromJSON(
		bytes,
		svc.GetServiceManager().GetEmptyInstanceDetails(),
		&pps,
	)
	instance.Service = svc
	instance.Plan = plan
	if instance.ParentAlias != "" {
		var parent service.Instance
		var ok bool
		parent, ok, err = s.GetInstanceByAlias(instance.ParentAlias)
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
	pipeline.SRem(instances, key)
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
	json, err := binding.ToJSON()
	if err != nil {
		return err
	}
	pipeline := s.redisClient.TxPipeline()
	pipeline.Set(key, json, 0)
	pipeline.SAdd(bindings, key)
	_, err = pipeline.Exec()
	if err != nil {
		return fmt.Errorf(
			`error writing binding "%s": %s`,
			binding.BindingID,
			err,
		)
	}
	return nil
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
	binding, err := service.NewBindingFromJSON(bytes, nil, nil)
	if err != nil {
		return binding, false, err
	}
	instance, ok, err := s.GetInstance(binding.InstanceID)
	if err != nil {
		return binding, false, err
	}
	// Now that we have schema for binding params, take a second pass at getting a
	// binding from the JSON
	if ok {
		bps := instance.Plan.GetSchemas().ServiceBindings.BindingParametersSchema
		binding, err = service.NewBindingFromJSON(
			bytes,
			instance.Service.GetServiceManager().GetEmptyBindingDetails(),
			&bps,
		)
	}
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

	pipeline := s.redisClient.TxPipeline()
	pipeline.Del(key)
	pipeline.SRem(bindings, key)
	_, err := pipeline.Exec()
	if err != nil {
		return false, fmt.Errorf(
			`error deleting binding "%s": %s`,
			bindingID,
			err,
		)
	}
	return true, nil
}

func getBindingKey(bindingID string) string {
	return fmt.Sprintf("bindings:%s", bindingID)
}

func (s *store) TestConnection() error {
	return s.redisClient.Ping().Err()
}
