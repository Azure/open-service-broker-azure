package storage

import (
	"log"
	"testing"

	"github.com/Azure/open-service-broker-azure/pkg/crypto/noop"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/fake"
	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var (
	noopCodec   = noop.NewCodec()
	redisClient = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	fakeServiceManager service.ServiceManager
	testStore          Store
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
	testStore = NewStore(
		redisClient,
		fakeCatalog,
		noopCodec,
	)
}

func TestWriteInstance(t *testing.T) {
	instance := getTestInstance()
	// First assert that the instance doesn't exist in Redis
	strCmd := redisClient.Get(instance.InstanceID)
	assert.Equal(t, redis.Nil, strCmd.Err())
	// Store the instance
	err := testStore.WriteInstance(instance)
	assert.Nil(t, err)
	// Assert that the instance is now in Redis
	strCmd = redisClient.Get(instance.InstanceID)
	assert.Nil(t, strCmd.Err())
}

func TestGetNonExistingInstance(t *testing.T) {
	instanceID := uuid.NewV4().String()
	// First assert that the instance doesn't exist in Redis
	strCmd := redisClient.Get(instanceID)
	assert.Equal(t, redis.Nil, strCmd.Err())
	// Try to retrieve the non-existing instance
	_, ok, err := testStore.GetInstance(instanceID)
	// Assert that the retrieval failed
	assert.False(t, ok)
	assert.Nil(t, err)
}

func TestGetExistingInstance(t *testing.T) {
	instance := getTestInstance()
	// First ensure the instance exists in Redis
	json, err := instance.ToJSON(noopCodec)
	assert.Nil(t, err)
	statCmd := redisClient.Set(instance.InstanceID, json, 0)
	assert.Nil(t, statCmd.Err())
	// Retrieve the instance
	retrievedInstance, ok, err := testStore.GetInstance(instance.InstanceID)
	// Assert that the retrieval was successful
	assert.Nil(t, err)
	assert.True(t, ok)
	// Blank out a few fields before we compare
	retrievedInstance.EncryptedProvisioningParameters = nil
	retrievedInstance.EncryptedUpdatingParameters = nil
	retrievedInstance.EncryptedProvisioningContext = nil
	assert.Equal(t, instance, retrievedInstance)
}

func TestDeleteNonExistingInstance(t *testing.T) {
	instanceID := uuid.NewV4().String()
	// First assert that the instance doesn't exist in Redis
	strCmd := redisClient.Get(instanceID)
	assert.Equal(t, redis.Nil, strCmd.Err())
	// Try to delete the non-existing instance
	ok, err := testStore.DeleteInstance(instanceID)
	// Assert that the delete failed
	assert.False(t, ok)
	assert.Nil(t, err)
}

func TestDeleteExistingInstance(t *testing.T) {
	instance := getTestInstance()
	// First ensure the instance exists in Redis
	json, err := instance.ToJSON(noopCodec)
	assert.Nil(t, err)
	statCmd := redisClient.Set(instance.InstanceID, json, 0)
	assert.Nil(t, statCmd.Err())
	// Delete the instance
	ok, err := testStore.DeleteInstance(instance.InstanceID)
	// Assert that the delete was successful
	assert.True(t, ok)
	assert.Nil(t, err)
	strCmd := redisClient.Get(instance.InstanceID)
	assert.Equal(t, redis.Nil, strCmd.Err())
}

func TestWriteBinding(t *testing.T) {
	binding := getTestBinding()
	// First assert that the binding doesn't exist in Redis
	strCmd := redisClient.Get(binding.BindingID)
	assert.Equal(t, redis.Nil, strCmd.Err())
	// Store the binding
	err := testStore.WriteBinding(binding)
	assert.Nil(t, err)
	// Assert that the binding is now in Redis
	strCmd = redisClient.Get(binding.BindingID)
	assert.Nil(t, strCmd.Err())
}

func TestGetNonExistingBinding(t *testing.T) {
	bindingID := uuid.NewV4().String()
	// First assert that the binding doesn't exist in Redis
	strCmd := redisClient.Get(bindingID)
	assert.Equal(t, redis.Nil, strCmd.Err())
	// Try to retrieve the non-existing binding
	_, ok, err := testStore.GetBinding(bindingID)
	// Assert that the retrieval failed
	assert.False(t, ok)
	assert.Nil(t, err)
}

func TestGetExistingBinding(t *testing.T) {
	binding := getTestBinding()
	// First ensure the binding exists in Redis
	json, err := binding.ToJSON(noopCodec)
	assert.Nil(t, err)
	statCmd := redisClient.Set(binding.BindingID, json, 0)
	assert.Nil(t, statCmd.Err())
	// Retrieve the binding
	retrievedBinding, ok, err := testStore.GetBinding(binding.BindingID)
	// Assert that the retrieval was successful
	assert.True(t, ok)
	assert.Nil(t, err)
	// Blank out a few fields before we compare
	retrievedBinding.EncryptedBindingParameters = nil
	retrievedBinding.EncryptedBindingContext = nil
	retrievedBinding.EncryptedCredentials = nil
	assert.Equal(t, binding, retrievedBinding)
}

func TestDeleteNonExistingBinding(t *testing.T) {
	bindingID := uuid.NewV4().String()
	// First assert that the binding doesn't exist in Redis
	strCmd := redisClient.Get(bindingID)
	assert.Equal(t, redis.Nil, strCmd.Err())
	// Try to delete the non-existing binding
	ok, err := testStore.DeleteBinding(bindingID)
	// Assert that the delete failed
	assert.False(t, ok)
	assert.Nil(t, err)
}

func TestDeleteExistingBinding(t *testing.T) {
	binding := getTestBinding()
	// First ensure the binding exists in Redis
	json, err := binding.ToJSON(noopCodec)
	assert.Nil(t, err)
	statCmd := redisClient.Set(binding.BindingID, json, 0)
	assert.Nil(t, statCmd.Err())
	// Delete the binding
	ok, err := testStore.DeleteBinding(binding.BindingID)
	// Assert that the delete was successful
	assert.True(t, ok)
	assert.Nil(t, err)
	strCmd := redisClient.Get(binding.BindingID)
	assert.Equal(t, redis.Nil, strCmd.Err())
}

func getTestInstance() service.Instance {
	return service.Instance{
		InstanceID: uuid.NewV4().String(),
		ServiceID:  fake.ServiceID,
		PlanID:     fake.StandardPlanID,
		StandardProvisioningParameters: service.StandardProvisioningParameters{
			Location:      "eastus",
			ResourceGroup: "test",
			Tags:          map[string]string{"foo": "bar"},
		},
		ProvisioningParameters: fakeServiceManager.GetEmptyProvisioningParameters(),
		UpdatingParameters:     fakeServiceManager.GetEmptyUpdatingParameters(),
		Status:                 service.InstanceStateProvisioned,
		StatusReason:           "",
		Location:               "eastus",
		ResourceGroup:          "test",
		Tags:                   map[string]string{"foo": "bar"},
		ProvisioningContext:    fakeServiceManager.GetEmptyProvisioningContext(),
	}
}

func getTestBinding() service.Binding {
	return service.Binding{
		BindingID:         uuid.NewV4().String(),
		InstanceID:        uuid.NewV4().String(),
		ServiceID:         fake.ServiceID,
		BindingParameters: fakeServiceManager.GetEmptyBindingParameters(),
		Status:            service.BindingStateBound,
		StatusReason:      "",
		BindingContext:    fakeServiceManager.GetEmptyBindingContext(),
		Credentials:       fakeServiceManager.GetEmptyCredentials(),
	}
}
