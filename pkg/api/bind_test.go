package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azure/azure-service-broker/pkg/crypto/noop"
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/services/echo"
	"github.com/stretchr/testify/assert"
)

func TestBindingWithInstanceThatDoesNotExist(t *testing.T) {
	s, err := getTestServer()
	assert.Nil(t, err)
	req, err := getBindingRequest(
		getDisposableInstanceID(),
		getDisposableBindingID(),
		nil,
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, responseEmptyJSON, rr.Body.Bytes())
}

func TestBindingWithInstanceThatIsNotFullyProvisioned(t *testing.T) {
	s, err := getTestServer()
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	planID := getDisposablePlanID()
	existingInstance := &service.Instance{
		InstanceID: instanceID,
		ServiceID:  getDisposableServiceID(),
		PlanID:     planID,
		Status:     service.InstanceStateProvisioning,
	}
	err = s.store.WriteInstance(existingInstance)
	assert.Nil(t, err)
	req, err := getBindingRequest(
		instanceID,
		getDisposableBindingID(),
		nil,
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	assert.Equal(t, responseEmptyJSON, rr.Body.Bytes())
}

func TestBindingWithServiceIDDifferentFromInstanceServiceID(t *testing.T) {
	s, err := getTestServer()
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	planID := getDisposablePlanID()
	existingInstance := &service.Instance{
		InstanceID: instanceID,
		ServiceID:  getDisposableServiceID(),
		PlanID:     planID,
		Status:     service.InstanceStateProvisioned,
	}
	err = s.store.WriteInstance(existingInstance)
	assert.Nil(t, err)
	req, err := getBindingRequest(
		instanceID,
		getDisposableBindingID(),
		&service.BindingRequest{
			ServiceID: getDisposableServiceID(),
			PlanID:    planID,
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Equal(t, responseEmptyJSON, rr.Body.Bytes())
}

func TestBindingWithPlanIDDifferentFromInstancePlanID(t *testing.T) {
	s, err := getTestServer()
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	serviceID := getDisposableServiceID()
	existingInstance := &service.Instance{
		InstanceID: instanceID,
		ServiceID:  serviceID,
		PlanID:     getDisposablePlanID(),
		Status:     service.InstanceStateProvisioned,
	}
	err = s.store.WriteInstance(existingInstance)
	assert.Nil(t, err)
	req, err := getBindingRequest(
		instanceID,
		getDisposableBindingID(),
		&service.BindingRequest{
			ServiceID: serviceID,
			PlanID:    getDisposablePlanID(),
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Equal(t, responseEmptyJSON, rr.Body.Bytes())
}

func TestBindingModuleNotFoundForServiceID(t *testing.T) {
	s, err := getTestServer()
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	serviceID := getDisposableServiceID()
	existingInstance := &service.Instance{
		InstanceID: instanceID,
		ServiceID:  serviceID,
		PlanID:     getDisposablePlanID(),
		Status:     service.InstanceStateProvisioned,
	}
	err = s.store.WriteInstance(existingInstance)
	assert.Nil(t, err)
	req, err := getBindingRequest(
		instanceID,
		getDisposableBindingID(),
		&service.BindingRequest{},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, responseEmptyJSON, rr.Body.Bytes())
}

func TestBindingWithExistingBindingWithDifferentInstanceID(
	t *testing.T,
) {
	s, err := getTestServer()
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	existingInstance := &service.Instance{
		InstanceID: instanceID,
		ServiceID:  echo.ServiceID,
		PlanID:     echo.StandardPlanID,
		Status:     service.InstanceStateProvisioned,
	}
	err = s.store.WriteInstance(existingInstance)
	assert.Nil(t, err)
	bindingID := getDisposableBindingID()
	existingBinding := &service.Binding{
		InstanceID: getDisposableInstanceID(),
		BindingID:  bindingID,
	}
	err = s.store.WriteBinding(existingBinding)
	assert.Nil(t, err)
	req, err := getBindingRequest(
		instanceID,
		bindingID,
		&service.BindingRequest{},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Equal(t, responseEmptyJSON, rr.Body.Bytes())
}

func TestBindingWithExistingBindingWithDifferentParameters(
	t *testing.T,
) {
	s, err := getTestServer()
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	existingInstance := &service.Instance{
		InstanceID: instanceID,
		ServiceID:  echo.ServiceID,
		PlanID:     echo.StandardPlanID,
		Status:     service.InstanceStateProvisioned,
	}
	err = s.store.WriteInstance(existingInstance)
	assert.Nil(t, err)
	bindingID := getDisposableBindingID()
	existingBinding := &service.Binding{
		InstanceID: instanceID,
		BindingID:  bindingID,
	}
	err = existingBinding.SetBindingParameters(
		&echo.BindingParameters{
			Message: "foo",
		},
		noop.NewCodec(),
	)
	assert.Nil(t, err)
	err = s.store.WriteBinding(existingBinding)
	assert.Nil(t, err)
	req, err := getBindingRequest(
		instanceID,
		bindingID,
		&service.BindingRequest{
			Parameters: &echo.BindingParameters{
				Message: "bar",
			},
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Equal(t, responseEmptyJSON, rr.Body.Bytes())
}

func TestBindingWithExistingBoundBindingWithSameAttributes(
	t *testing.T,
) {
	s, err := getTestServer()
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	existingInstance := &service.Instance{
		InstanceID: instanceID,
		ServiceID:  echo.ServiceID,
		PlanID:     echo.StandardPlanID,
		Status:     service.InstanceStateProvisioned,
	}
	err = s.store.WriteInstance(existingInstance)
	assert.Nil(t, err)
	bindingID := getDisposableBindingID()
	existingBinding := &service.Binding{
		InstanceID: instanceID,
		BindingID:  bindingID,
		Status:     service.BindingStateBound,
	}
	err = s.store.WriteBinding(existingBinding)
	assert.Nil(t, err)
	req, err := getBindingRequest(
		instanceID,
		bindingID,
		&service.BindingRequest{},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	// TODO: Test the response body
}

func TestBindingWithExistingFailedBindingWithSameAttributes(
	t *testing.T,
) {
	s, err := getTestServer()
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	existingInstance := &service.Instance{
		InstanceID: instanceID,
		ServiceID:  echo.ServiceID,
		PlanID:     echo.StandardPlanID,
		Status:     service.InstanceStateProvisioned,
	}
	err = s.store.WriteInstance(existingInstance)
	assert.Nil(t, err)
	bindingID := getDisposableBindingID()
	existingBinding := &service.Binding{
		InstanceID: instanceID,
		BindingID:  bindingID,
		Status:     service.BindingStateBindingFailed,
	}
	err = s.store.WriteBinding(existingBinding)
	assert.Nil(t, err)
	req, err := getBindingRequest(
		instanceID,
		bindingID,
		&service.BindingRequest{},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Equal(t, responseEmptyJSON, rr.Body.Bytes())
}

func TestBrandNewBinding(t *testing.T) {
	s, err := getTestServer()
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	existingInstance := &service.Instance{
		InstanceID: instanceID,
		ServiceID:  echo.ServiceID,
		PlanID:     echo.StandardPlanID,
		Status:     service.InstanceStateProvisioned,
	}
	err = s.store.WriteInstance(existingInstance)
	assert.Nil(t, err)
	req, err := getBindingRequest(
		instanceID,
		getDisposableBindingID(),
		&service.BindingRequest{},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)
	// TODO: Test the response body
}

func getBindingRequest(
	instanceID string,
	bindingID string,
	br *service.BindingRequest,
) (*http.Request, error) {
	bodyStr := ""
	if br != nil {
		var err error
		bodyStr, err = br.ToJSONString()
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf(
			"/v2/service_instances/%s/service_bindings/%s",
			instanceID,
			bindingID,
		),
		bytes.NewBuffer([]byte(bodyStr)),
	)
	if err != nil {
		return nil, err
	}
	return req, nil
}
