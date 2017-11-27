package storage

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var (
	redisClient = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	testStore = NewStore(redisClient)
)

func TestWriteInstance(t *testing.T) {
	instanceID := getDisposableInstanceID()
	// First assert that the instance doesn't exist in Redis
	strCmd := redisClient.Get(instanceID)
	assert.Equal(t, redis.Nil, strCmd.Err())
	// Store the instance
	err := testStore.WriteInstance(&service.Instance{
		InstanceID: instanceID,
	})
	assert.Nil(t, err)
	// Assert that the instance is now in Redis
	strCmd = redisClient.Get(instanceID)
	assert.Nil(t, strCmd.Err())
}

func TestGetNonExistingInstance(t *testing.T) {
	instanceID := getDisposableInstanceID()
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
	instanceID := getDisposableInstanceID()
	// First ensure the instance exists in Redis
	statCmd := redisClient.Set(instanceID, getInstanceJSON(instanceID), 0)
	assert.Nil(t, statCmd.Err())
	// Retrieve the instance
	instance, ok, err := testStore.GetInstance(instanceID)
	// Assert that the retrieval was successful
	assert.True(t, ok)
	assert.Nil(t, err)
	// Asset that instance is not nil before using
	if assert.NotNil(t, instance, "instance should not be nil") {
		assert.Equal(t, instanceID, instance.InstanceID)
	}
}

func TestDeleteNonExistingInstance(t *testing.T) {
	instanceID := getDisposableInstanceID()
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
	instanceID := getDisposableInstanceID()
	// First ensure the instance exists in Redis
	statCmd := redisClient.Set(instanceID, getInstanceJSON(instanceID), 0)
	assert.Nil(t, statCmd.Err())
	// Delete the instance
	ok, err := testStore.DeleteInstance(instanceID)
	// Assert that the delete was successful
	assert.True(t, ok)
	assert.Nil(t, err)
	strCmd := redisClient.Get(instanceID)
	assert.Equal(t, redis.Nil, strCmd.Err())
}

func TestWriteBinding(t *testing.T) {
	bindingID := getDisposableBindingID()
	// First assert that the binding doesn't exist in Redis
	strCmd := redisClient.Get(bindingID)
	assert.Equal(t, redis.Nil, strCmd.Err())
	// Store the binding
	err := testStore.WriteBinding(&service.Binding{
		BindingID: bindingID,
	})
	assert.Nil(t, err)
	// Assert that the binding is now in Redis
	strCmd = redisClient.Get(bindingID)
	assert.Nil(t, strCmd.Err())
}

func TestGetNonExistingBinding(t *testing.T) {
	bindingID := getDisposableBindingID()
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
	bindingID := getDisposableBindingID()
	// First ensure the binding exists in Redis
	statCmd := redisClient.Set(bindingID, getBindingJSON(bindingID), 0)
	assert.Nil(t, statCmd.Err())
	// Retrieve the binding
	binding, ok, err := testStore.GetBinding(bindingID)
	// Assert that the retrieval was successful
	assert.True(t, ok)
	assert.Nil(t, err)
	// Assert that binding is not nil before using
	if assert.NotNil(t, binding, "binding should not be nil") {
		assert.Equal(t, bindingID, binding.BindingID)
	}
}

func TestDeleteNonExistingBinding(t *testing.T) {
	bindingID := getDisposableBindingID()
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
	bindingID := getDisposableBindingID()
	// First ensure the binding exists in Redis
	statCmd := redisClient.Set(bindingID, getBindingJSON(bindingID), 0)
	assert.Nil(t, statCmd.Err())
	// Delete the binding
	ok, err := testStore.DeleteBinding(bindingID)
	// Assert that the delete was successful
	assert.True(t, ok)
	assert.Nil(t, err)
	strCmd := redisClient.Get(bindingID)
	assert.Equal(t, redis.Nil, strCmd.Err())
}

func getInstanceJSON(instanceID string) string {
	return fmt.Sprintf(`{"instanceId":"%s"}`, instanceID)
}

func getBindingJSON(bindingID string) string {
	return fmt.Sprintf(`{"bindingId":"%s"}`, bindingID)
}

func getDisposableInstanceID() string {
	return uuid.NewV4().String()
}

func getDisposableBindingID() string {
	return uuid.NewV4().String()
}
