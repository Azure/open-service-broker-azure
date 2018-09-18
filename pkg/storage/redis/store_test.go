package redis

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/fake"
	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var (
	fakeServiceManager service.ServiceManager
	testStore          *store
	config             Config
)

func init() {
	var err error
	fakeModule, err := fake.New()
	if err != nil {
		log.Fatal(err)
	}
	fakeCatalog, err := fakeModule.GetCatalog()
	if err != nil {
		log.Fatal(err)
	}
	fakeServiceManager = fakeModule.ServiceManager
	config = NewConfigWithDefaults()
	config.RedisHost = os.Getenv("STORAGE_REDIS_HOST")
	config.RedisPrefix = uuid.NewV4().String()
	str, err := NewStore(
		fakeCatalog,
		config,
	)
	if err != nil {
		log.Fatal(err)
	}
	testStore = str.(*store)
}

func TestWriteInstance(t *testing.T) {
	instance := getTestInstance()
	key := testStore.getInstanceKey(instance.InstanceID)
	// First assert that the instance doesn't exist in Redis
	strCmd := testStore.redisClient.Get(key)
	assert.Equal(t, redis.Nil, strCmd.Err())
	// Store the instance
	err := testStore.WriteInstance(instance)
	assert.Nil(t, err)
	// Assert that the instance is now in Redis
	strCmd = testStore.redisClient.Get(key)
	assert.Nil(t, strCmd.Err())
	// TODO: krancour: This next assertion only holds true if the Redis DB is
	// reset between test runs. Need to fix this at some point.
	// sCmd := testStore.redisClient.SMembers(instances)
	// assert.Nil(t, strCmd.Err())
	// count, err := sCmd.Result()
	// assert.Nil(t, err)
	// assert.Equal(t, 1, len(count))
	boolCmd := testStore.redisClient.SIsMember(
		wrapKey(
			config.RedisPrefix,
			"instances",
		),
		key,
	)
	assert.Nil(t, boolCmd.Err())
	found, _ := boolCmd.Result()
	assert.True(t, found)
}

func TestWriteInstanceWithAlias(t *testing.T) {
	instance := getTestInstance()
	instance.Alias = uuid.NewV4().String()
	key := testStore.getInstanceKey(instance.InstanceID)
	// First assert that the instance doesn't exist in Redis
	strCmd := testStore.redisClient.Get(key)
	assert.Equal(t, redis.Nil, strCmd.Err())
	// Nor does its alias
	aliasKey := testStore.getInstanceAliasKey(instance.Alias)
	strCmd = testStore.redisClient.Get(aliasKey)
	assert.Equal(t, redis.Nil, strCmd.Err())
	// Store the instance
	err := testStore.WriteInstance(instance)
	assert.Nil(t, err)
	// Assert that the instance is now in Redis
	strCmd = testStore.redisClient.Get(key)
	assert.Nil(t, strCmd.Err())
	// Assert that the alias is as well
	strCmd = testStore.redisClient.Get(aliasKey)
	assert.Nil(t, strCmd.Err())
	instanceID, err := strCmd.Result()
	assert.Nil(t, err)
	assert.Equal(t, instance.InstanceID, instanceID)
}

func TestWriteInstanceWithParent(t *testing.T) {
	instance := getTestInstance()
	instance.ParentAlias = uuid.NewV4().String()
	key := testStore.getInstanceKey(instance.InstanceID)
	// First assert that the instance doesn't exist in Redis
	strCmd := testStore.redisClient.Get(key)
	assert.Equal(t, redis.Nil, strCmd.Err())
	// Nor does any index of parent alias to children
	parentAliasChildrenKey := testStore.getInstanceAliasChildrenKey(
		instance.ParentAlias,
	)
	boolCmd :=
		testStore.redisClient.SIsMember(parentAliasChildrenKey, instance.InstanceID)
	assert.Nil(t, boolCmd.Err())
	childFoundInIndex, err := boolCmd.Result()
	assert.Nil(t, err)
	assert.False(t, childFoundInIndex)
	// Store the instance
	err = testStore.WriteInstance(instance)
	assert.Nil(t, err)
	// Assert that the instance is now in Redis
	strCmd = testStore.redisClient.Get(key)
	assert.Nil(t, strCmd.Err())
	// And the index for parent alias to children contains this instance
	boolCmd =
		testStore.redisClient.SIsMember(parentAliasChildrenKey, instance.InstanceID)
	assert.Nil(t, boolCmd.Err())
	childFoundInIndex, err = boolCmd.Result()
	assert.Nil(t, err)
	assert.True(t, childFoundInIndex)
}

func TestGetNonExistingInstance(t *testing.T) {
	instanceID := uuid.NewV4().String()
	key := testStore.getInstanceKey(instanceID)
	// First assert that the instance doesn't exist in Redis
	strCmd := testStore.redisClient.Get(key)
	assert.Equal(t, redis.Nil, strCmd.Err())
	// Try to retrieve the non-existing instance
	_, ok, err := testStore.GetInstance(instanceID)
	// Assert that the retrieval failed
	assert.False(t, ok)
	assert.Nil(t, err)
}

func TestGetExistingInstance(t *testing.T) {
	instance := getTestInstance()
	key := testStore.getInstanceKey(instance.InstanceID)
	// First ensure the instance exists in Redis
	json, err := instance.ToJSON()
	assert.Nil(t, err)
	statCmd := testStore.redisClient.Set(key, json, 0)
	assert.Nil(t, statCmd.Err())
	// Retrieve the instance
	retrievedInstance, ok, err := testStore.GetInstance(instance.InstanceID)
	// Assert that the retrieval was successful
	assert.Nil(t, err)
	assert.True(t, ok)
	// Blank out a few fields before we compare
	retrievedInstance.Service = nil
	retrievedInstance.Plan = nil
	assert.Equal(t, instance, retrievedInstance)
}

func TestGetExistingInstanceWithParent(t *testing.T) {
	// Make a parent instance
	parentInstance := getTestInstance()
	parentInstance.Alias = uuid.NewV4().String()
	parentKey := testStore.getInstanceKey(parentInstance.InstanceID)
	// Ensure the parent instance exists in Redis
	json, err := parentInstance.ToJSON()
	assert.Nil(t, err)
	statCmd := testStore.redisClient.Set(parentKey, json, 0)
	assert.Nil(t, statCmd.Err())
	// Ensure the parent instance's alias also exists in Redis
	parentAliasKey := testStore.getInstanceAliasKey(parentInstance.Alias)
	statCmd =
		testStore.redisClient.Set(parentAliasKey, parentInstance.InstanceID, 0)
	assert.Nil(t, statCmd.Err())
	// Make a child instance
	instance := getTestInstance()
	instance.ParentAlias = parentInstance.Alias
	instance.Parent = &parentInstance
	key := testStore.getInstanceKey(instance.InstanceID)
	// Ensure the child instance exists in Redis
	json, err = instance.ToJSON()
	assert.Nil(t, err)
	statCmd = testStore.redisClient.Set(key, json, 0)
	assert.Nil(t, statCmd.Err())
	// Retrieve the child instance
	retrievedInstance, ok, err := testStore.GetInstance(instance.InstanceID)
	// Assert that the retrieval was successful
	assert.Nil(t, err)
	assert.True(t, ok)
	// Blank out a few fields before we compare
	retrievedInstance.Service = nil
	retrievedInstance.Parent.Service = nil
	retrievedInstance.Plan = nil
	retrievedInstance.Parent.Plan = nil
	assert.Equal(t, instance, retrievedInstance)
}

func TestGetNonExistingInstanceByAlias(t *testing.T) {
	alias := uuid.NewV4().String()
	aliasKey := testStore.getInstanceAliasKey(alias)
	// First assert that the alias doesn't exist in Redis
	strCmd := testStore.redisClient.Get(aliasKey)
	assert.Equal(t, redis.Nil, strCmd.Err())
	// Try to retrieve the non-existing instance by alias
	_, ok, err := testStore.GetInstanceByAlias(aliasKey)
	// Assert that the retrieval failed
	assert.False(t, ok)
	assert.Nil(t, err)
}

func TestGetExistingInstanceByAlias(t *testing.T) {
	instance := getTestInstance()
	key := testStore.getInstanceKey(instance.InstanceID)
	// First ensure the instance exists in Redis
	json, err := instance.ToJSON()
	assert.Nil(t, err)
	statCmd := testStore.redisClient.Set(key, json, 0)
	assert.Nil(t, statCmd.Err())
	// And so does the alias
	aliasKey := testStore.getInstanceAliasKey(instance.Alias)
	statCmd = testStore.redisClient.Set(aliasKey, instance.InstanceID, 0)
	assert.Nil(t, statCmd.Err())
	// Retrieve the instance by alias
	retrievedInstance, ok, err := testStore.GetInstanceByAlias(instance.Alias)
	// Assert that the retrieval was successful
	assert.Nil(t, err)
	assert.True(t, ok)
	// Blank out a few fields before we compare
	retrievedInstance.Service = nil
	retrievedInstance.Plan = nil
	assert.Equal(t, instance, retrievedInstance)
}

func TestDeleteNonExistingInstance(t *testing.T) {
	instanceID := uuid.NewV4().String()
	key := testStore.getInstanceKey(instanceID)
	// First assert that the instance doesn't exist in Redis
	strCmd := testStore.redisClient.Get(key)
	assert.Equal(t, redis.Nil, strCmd.Err())
	// Try to delete the non-existing instance
	ok, err := testStore.DeleteInstance(instanceID)
	// Assert that the delete failed
	assert.False(t, ok)
	assert.Nil(t, err)
}

func TestDeleteExistingInstance(t *testing.T) {
	instance := getTestInstance()
	key := testStore.getInstanceKey(instance.InstanceID)
	// First ensure the instance exists in Redis
	json, err := instance.ToJSON()
	assert.Nil(t, err)
	statCmd := testStore.redisClient.Set(key, json, 0)
	assert.Nil(t, statCmd.Err())
	// Delete the instance
	ok, err := testStore.DeleteInstance(instance.InstanceID)
	// Assert that the delete was successful
	assert.True(t, ok)
	assert.Nil(t, err)
	strCmd := testStore.redisClient.Get(key)
	assert.Equal(t, redis.Nil, strCmd.Err())
	boolCmd := testStore.redisClient.SIsMember(
		wrapKey(
			config.RedisPrefix,
			"instances",
		),
		key,
	)
	assert.Nil(t, boolCmd.Err())
	found, _ := boolCmd.Result()
	assert.False(t, found)
}

func TestDeleteExistingInstanceWithAlias(t *testing.T) {
	instance := getTestInstance()
	instance.Alias = uuid.NewV4().String()
	key := testStore.getInstanceKey(instance.InstanceID)
	// First ensure the instance exists in Redis
	json, err := instance.ToJSON()
	assert.Nil(t, err)
	statCmd := testStore.redisClient.Set(key, json, 0)
	assert.Nil(t, statCmd.Err())
	// And so does the alias
	aliasKey := testStore.getInstanceAliasKey(instance.Alias)
	statCmd = testStore.redisClient.Set(aliasKey, instance.InstanceID, 0)
	assert.Nil(t, statCmd.Err())
	// Delete the instance
	ok, err := testStore.DeleteInstance(instance.InstanceID)
	// Assert that the delete was successful
	assert.True(t, ok)
	assert.Nil(t, err)
	strCmd := testStore.redisClient.Get(key)
	assert.Equal(t, redis.Nil, strCmd.Err())
	// Assert that the alias is also gone
	strCmd = testStore.redisClient.Get(aliasKey)
	assert.Equal(t, redis.Nil, strCmd.Err())
}

func TestDeleteExistingInstanceWithParent(t *testing.T) {
	// Make a parent instance
	parentInstance := getTestInstance()
	parentKey := testStore.getInstanceKey(parentInstance.InstanceID)
	// Ensure the parent instance exists in Redis
	json, err := parentInstance.ToJSON()
	assert.Nil(t, err)
	statCmd := testStore.redisClient.Set(parentKey, json, 0)
	assert.Nil(t, statCmd.Err())
	// Ensure the parent instance's alias also exists in Redis
	parentAliasKey := testStore.getInstanceAliasKey(parentInstance.Alias)
	statCmd =
		testStore.redisClient.Set(parentAliasKey, parentInstance.InstanceID, 0)
	assert.Nil(t, statCmd.Err())
	// Make a child instance
	instance := getTestInstance()
	instance.ParentAlias = parentInstance.Alias
	instance.Parent = &parentInstance
	key := testStore.getInstanceKey(instance.InstanceID)
	// Ensure the child instance exists in Redis
	json, err = instance.ToJSON()
	assert.Nil(t, err)
	statCmd = testStore.redisClient.Set(key, json, 0)
	assert.Nil(t, statCmd.Err())
	// Delete the instance
	ok, err := testStore.DeleteInstance(instance.InstanceID)
	// Assert that the delete was successful
	assert.True(t, ok)
	assert.Nil(t, err)
	strCmd := testStore.redisClient.Get(key)
	assert.Equal(t, redis.Nil, strCmd.Err())
	// And the index of parent alias to children no longer contains this instance
	parentAliasChildrenKey := testStore.getInstanceAliasChildrenKey(
		instance.ParentAlias,
	)
	boolCmd :=
		testStore.redisClient.SIsMember(parentAliasChildrenKey, instance.InstanceID)
	assert.Nil(t, boolCmd.Err())
	childFoundInIndex, err := boolCmd.Result()
	assert.Nil(t, err)
	assert.False(t, childFoundInIndex)
}

func TestGetInstanceChildCountByAlias(t *testing.T) {
	const count = 5
	instanceAlias := uuid.NewV4().String()
	instanceAliasChildrenKey := testStore.getInstanceAliasChildrenKey(
		instanceAlias,
	)
	for i := 0; i < count; i++ {
		// Add a new, unique, child instance ID to the index
		testStore.redisClient.SAdd(instanceAliasChildrenKey, uuid.NewV4().String())
		// Count the children
		children, err := testStore.GetInstanceChildCountByAlias(instanceAlias)
		assert.Nil(t, err)
		// Assert the size of the index is what we expect
		assert.Equal(t, int64(i+1), children)
	}
}

func TestWriteBinding(t *testing.T) {
	binding := getTestBinding()
	key := testStore.getBindingKey(binding.BindingID)
	// First assert that the binding doesn't exist in Redis
	strCmd := testStore.redisClient.Get(key)
	assert.Equal(t, redis.Nil, strCmd.Err())
	// Store the binding
	err := testStore.WriteBinding(binding)
	assert.Nil(t, err)
	// Assert that the binding is now in Redis
	strCmd = testStore.redisClient.Get(key)
	assert.Nil(t, strCmd.Err())
	boolCmd := testStore.redisClient.SIsMember(
		wrapKey(
			config.RedisPrefix,
			"bindings",
		),
		key,
	)
	assert.Nil(t, boolCmd.Err())
	found, _ := boolCmd.Result()
	assert.True(t, found)
}

func TestGetNonExistingBinding(t *testing.T) {
	bindingID := uuid.NewV4().String()
	key := testStore.getBindingKey(bindingID)
	// First assert that the binding doesn't exist in Redis
	strCmd := testStore.redisClient.Get(key)
	assert.Equal(t, redis.Nil, strCmd.Err())
	// Try to retrieve the non-existing binding
	_, ok, err := testStore.GetBinding(bindingID)
	// Assert that the retrieval failed
	assert.False(t, ok)
	assert.Nil(t, err)
}

func TestGetExistingBinding(t *testing.T) {
	binding := getTestBinding()
	key := testStore.getBindingKey(binding.BindingID)
	// First ensure the binding exists in Redis
	json, err := binding.ToJSON()
	assert.Nil(t, err)

	statCmd := testStore.redisClient.Set(key, json, 0)
	assert.Nil(t, statCmd.Err())
	// Retrieve the binding
	retrievedBinding, ok, err := testStore.GetBinding(binding.BindingID)
	// Assert that the retrieval was successful
	assert.True(t, ok)
	assert.Nil(t, err)
	// Blank out a few fields before we compare
	assert.Equal(t, binding, retrievedBinding)
}

func TestDeleteNonExistingBinding(t *testing.T) {
	bindingID := uuid.NewV4().String()
	key := testStore.getBindingKey(bindingID)
	// First assert that the binding doesn't exist in Redis
	strCmd := testStore.redisClient.Get(key)
	assert.Equal(t, redis.Nil, strCmd.Err())
	// Try to delete the non-existing binding
	ok, err := testStore.DeleteBinding(bindingID)
	// Assert that the delete failed
	assert.False(t, ok)
	assert.Nil(t, err)
}

func TestDeleteExistingBinding(t *testing.T) {
	binding := getTestBinding()
	key := testStore.getBindingKey(binding.BindingID)
	// First ensure the binding exists in Redis
	json, err := binding.ToJSON()
	assert.Nil(t, err)
	statCmd := testStore.redisClient.Set(key, json, 0)
	assert.Nil(t, statCmd.Err())
	// Delete the binding
	ok, err := testStore.DeleteBinding(binding.BindingID)
	// Assert that the delete was successful
	assert.True(t, ok)
	assert.Nil(t, err)
	strCmd := testStore.redisClient.Get(key)
	assert.Equal(t, redis.Nil, strCmd.Err())
}

func TestGetInstanceKey(t *testing.T) {
	const rawKey = "foo"
	expected := fmt.Sprintf("%s:instances:%s", config.RedisPrefix, rawKey)
	assert.Equal(t, expected, testStore.getInstanceKey(rawKey))
}

func TestGetBindingKey(t *testing.T) {
	const rawKey = "foo"
	expected := fmt.Sprintf("%s:bindings:%s", config.RedisPrefix, rawKey)
	assert.Equal(t, expected, testStore.getBindingKey(rawKey))
}

func getTestInstance() service.Instance {
	return service.Instance{
		InstanceID:   uuid.NewV4().String(),
		ServiceID:    fake.ServiceID,
		PlanID:       fake.StandardPlanID,
		Status:       service.InstanceStateProvisioned,
		StatusReason: "",
	}
}

func getTestBinding() service.Binding {
	return service.Binding{
		BindingID:    uuid.NewV4().String(),
		InstanceID:   uuid.NewV4().String(),
		ServiceID:    fake.ServiceID,
		Status:       service.BindingStateBound,
		StatusReason: "",
	}
}
