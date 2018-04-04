package redis

import (
	"crypto/tls"
	"errors"
	"fmt"

	"github.com/Azure/open-service-broker-azure/pkg/crypto"
	"github.com/Azure/open-service-broker-azure/pkg/crypto/aes256"
	"github.com/Azure/open-service-broker-azure/pkg/crypto/noop"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/storage"
	log "github.com/Sirupsen/logrus"
	"github.com/go-redis/redis"
)

type store struct {
	redisClient *redis.Client
	catalog     service.Catalog
	codec       crypto.Codec
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

	var codec crypto.Codec
	switch config.EncryptionScheme {
	case crypto.AES256:
		if config.AES256Key == "" {
			return nil, errors.New("AES256 key was not specified")
		}
		if len(config.AES256Key) != 32 {
			return nil, errors.New("AES256 key is an invalid length")
		}
		var err error
		codec, err = aes256.NewCodec([]byte(config.AES256Key))
		if err != nil {
			return nil, err
		}
		log.WithField(
			"encryptionScheme",
			config.EncryptionScheme,
		).Info("Sensitive instance and binding details will be encrypted")
	case crypto.NOOP:
		codec = noop.NewCodec()
		log.Warn(
			"ENCRYPTION IS DISABLED -- THIS IS NOT A SUITABLE OPTION FOR PRODUCTION",
		)
	default:
		return nil, fmt.Errorf(
			`unrecognized encryption scheme "%s"`,
			config.EncryptionScheme,
		)
	}

	return &store{
		redisClient: redis.NewClient(redisOpts),
		catalog:     catalog,
		codec:       codec,
	}, nil
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
