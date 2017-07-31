package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	fakeAsync "github.com/Azure/azure-service-broker/pkg/async/fake"
	"github.com/Azure/azure-service-broker/pkg/crypto/noop"
	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/Azure/azure-service-broker/pkg/services/echo"
	"github.com/stretchr/testify/assert"
)

func TestProvisioningWithAcceptIncompleteNotSet(t *testing.T) {
	s, err := getTestServer()
	assert.Nil(t, err)
	req, err := getProvisionRequest(getDisposableInstanceID(), nil, nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	assert.Equal(t, responseAsyncRequired, rr.Body.Bytes())
}

func TestProvisioningWithAcceptIncompleteNotTrue(t *testing.T) {
	s, err := getTestServer()
	assert.Nil(t, err)
	req, err := getProvisionRequest(
		getDisposableInstanceID(),
		map[string]string{
			"accepts_incomplete": "false",
		},
		nil,
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	assert.Equal(t, responseAsyncRequired, rr.Body.Bytes())
}

func TestProvisioningWithMissingServiceID(t *testing.T) {
	s, err := getTestServer()
	assert.Nil(t, err)
	req, err := getProvisionRequest(
		getDisposableInstanceID(),
		map[string]string{
			"accepts_incomplete": "true",
		},
		&service.ProvisioningRequest{
			ServiceID: "",
			PlanID:    getDisposablePlanID(),
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, responseServiceIDRequired, rr.Body.Bytes())
}

func TestProvisioningWithMissingPlanID(t *testing.T) {
	s, err := getTestServer()
	assert.Nil(t, err)
	req, err := getProvisionRequest(
		getDisposableInstanceID(),
		map[string]string{
			"accepts_incomplete": "true",
		},
		&service.ProvisioningRequest{
			ServiceID: getDisposableServiceID(),
			PlanID:    "",
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, responsePlanIDRequired, rr.Body.Bytes())
}

func TestProvisioningWithInvalidServiceID(t *testing.T) {
	s, err := getTestServer()
	assert.Nil(t, err)
	req, err := getProvisionRequest(
		getDisposableInstanceID(),
		map[string]string{
			"accepts_incomplete": "true",
		},
		&service.ProvisioningRequest{
			ServiceID: getDisposableServiceID(),
			PlanID:    getDisposablePlanID(),
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, responseInvalidServiceID, rr.Body.Bytes())
}

func TestProvisioningWithInvalidPlanID(t *testing.T) {
	s, err := getTestServer()
	assert.Nil(t, err)
	req, err := getProvisionRequest(
		getDisposableInstanceID(),
		map[string]string{
			"accepts_incomplete": "true",
		},
		&service.ProvisioningRequest{
			ServiceID: echo.ServiceID,
			PlanID:    getDisposablePlanID(),
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, responseInvalidPlanID, rr.Body.Bytes())
}

func TestProvisioningWithExistingInstanceWithDifferentAttributes(
	t *testing.T,
) {
	s, err := getTestServer()
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	existingInstance := &service.Instance{
		InstanceID: instanceID,
		ServiceID:  echo.ServiceID,
		PlanID:     echo.StandardPlanID,
	}
	err = existingInstance.SetProvisioningParameters(
		&echo.ProvisioningParameters{
			Message: "foo",
		},
		noop.NewCodec(),
	)
	assert.Nil(t, err)
	err = s.store.WriteInstance(existingInstance)
	assert.Nil(t, err)
	req, err := getProvisionRequest(
		instanceID,
		map[string]string{
			"accepts_incomplete": "true",
		},
		&service.ProvisioningRequest{
			ServiceID: echo.ServiceID,
			PlanID:    echo.StandardPlanID,
			Parameters: &echo.ProvisioningParameters{
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

func TestProvisioningWithExistingInstanceWithSameAttributesAndFullyProvisioned(
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
	req, err := getProvisionRequest(
		instanceID,
		map[string]string{
			"accepts_incomplete": "true",
		},
		&service.ProvisioningRequest{
			ServiceID: echo.ServiceID,
			PlanID:    echo.StandardPlanID,
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, responseEmptyJSON, rr.Body.Bytes())
}

func TestProvisioningWithExistingInstanceWithSameAttributesAndNotFullyProvisioned( // nolint: lll
	t *testing.T,
) {
	s, err := getTestServer()
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	existingInstance := &service.Instance{
		InstanceID: instanceID,
		ServiceID:  echo.ServiceID,
		PlanID:     echo.StandardPlanID,
		Status:     service.InstanceStateProvisioning,
	}
	err = s.store.WriteInstance(existingInstance)
	assert.Nil(t, err)
	req, err := getProvisionRequest(
		instanceID,
		map[string]string{
			"accepts_incomplete": "true",
		},
		&service.ProvisioningRequest{
			ServiceID: echo.ServiceID,
			PlanID:    echo.StandardPlanID,
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusAccepted, rr.Code)
	assert.Equal(t, responseEmptyJSON, rr.Body.Bytes())
}

func TestKickOffNewAsyncProvisioning(t *testing.T) {
	s, err := getTestServer()
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	req, err := getProvisionRequest(
		instanceID,
		map[string]string{
			"accepts_incomplete": "true",
		},
		&service.ProvisioningRequest{
			ServiceID: echo.ServiceID,
			PlanID:    echo.StandardPlanID,
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusAccepted, rr.Code)
	e := s.asyncEngine.(*fakeAsync.Engine)
	assert.Equal(t, 1, len(e.SubmittedTasks))
	assert.Equal(t, responseEmptyJSON, rr.Body.Bytes())
}

func getProvisionRequest(
	instanceID string,
	queryParams map[string]string,
	pr *service.ProvisioningRequest,
) (*http.Request, error) {
	var body []byte
	if pr != nil {
		var err error
		body, err = pr.ToJSON()
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/v2/service_instances/%s", instanceID),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}
	if queryParams != nil {
		q := req.URL.Query()
		for k, v := range queryParams {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
	return req, nil
}
