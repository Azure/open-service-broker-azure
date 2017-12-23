package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azure/open-service-broker-azure/pkg/services/fake"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/stretchr/testify/assert"
)

func TestUnbindingBindingThatIsNotFound(t *testing.T) {
	s, _, err := getTestServer("", "")
	assert.Nil(t, err)
	req, err := getUnbindingRequest(
		getDisposableInstanceID(),
		getDisposableBindingID(),
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusGone, rr.Code)
	assert.Equal(t, responseEmptyJSON, rr.Body.Bytes())
}

func TestUnbindingWithInstanceIDDifferentFromBindingInstanceID(t *testing.T) {
	s, _, err := getTestServer("", "")
	assert.Nil(t, err)
	bindingID := getDisposableBindingID()
	err = s.store.WriteBinding(service.Binding{
		InstanceID: getDisposableInstanceID(),
		BindingID:  bindingID,
		ServiceID:  fake.ServiceID,
	})
	assert.Nil(t, err)
	req, err := getUnbindingRequest(
		getDisposableInstanceID(),
		bindingID,
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Equal(t, responseEmptyJSON, rr.Body.Bytes())
}

func TestUnbindingFromInstanceThatDoesNotExist(t *testing.T) {
	s, _, err := getTestServer("", "")
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	bindingID := getDisposableBindingID()
	err = s.store.WriteBinding(service.Binding{
		InstanceID: instanceID,
		BindingID:  bindingID,
		ServiceID:  fake.ServiceID,
	})
	assert.Nil(t, err)
	req, err := getUnbindingRequest(
		instanceID,
		bindingID,
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, responseEmptyJSON, rr.Body.Bytes())
	_, ok, err := s.store.GetBinding(bindingID)
	assert.Nil(t, err)
	assert.False(t, ok)
}

func TestUnbindingFromInstanceThatExists(t *testing.T) {
	s, m, err := getTestServer("", "")
	assert.Nil(t, err)
	unbindCalled := false
	m.ServiceManager.UnbindBehavior = func(
		service.Instance,
		service.BindingDetails,
	) error {
		unbindCalled = true
		return nil
	}
	instanceID := getDisposableInstanceID()
	bindingID := getDisposableBindingID()
	err = s.store.WriteInstance(service.Instance{
		InstanceID: instanceID,
		ServiceID:  fake.ServiceID,
		PlanID:     fake.StandardPlanID,
	})
	assert.Nil(t, err)
	err = s.store.WriteBinding(service.Binding{
		InstanceID: instanceID,
		BindingID:  bindingID,
		ServiceID:  fake.ServiceID,
	})
	assert.Nil(t, err)
	req, err := getUnbindingRequest(
		instanceID,
		bindingID,
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, responseEmptyJSON, rr.Body.Bytes())
	assert.True(t, unbindCalled)
	_, ok, err := s.store.GetBinding(bindingID)
	assert.Nil(t, err)
	assert.False(t, ok)
}

func getUnbindingRequest(
	instanceID string,
	bindingID string,
) (*http.Request, error) {
	return http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf(
			"/v2/service_instances/%s/service_bindings/%s",
			instanceID,
			bindingID,
		),
		nil,
	)
}
