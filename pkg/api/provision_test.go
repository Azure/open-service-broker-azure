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
	"github.com/Azure/azure-service-broker/pkg/services/fake"
	"github.com/stretchr/testify/assert"
)

func TestProvisioningWithAcceptIncompleteNotSet(t *testing.T) {
	s, _, err := getTestServer()
	assert.Nil(t, err)
	req, err := getProvisionRequest(getDisposableInstanceID(), nil, nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	assert.Equal(t, responseAsyncRequired, rr.Body.Bytes())
}

func TestProvisioningWithAcceptIncompleteNotTrue(t *testing.T) {
	s, _, err := getTestServer()
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
	s, _, err := getTestServer()
	assert.Nil(t, err)
	req, err := getProvisionRequest(
		getDisposableInstanceID(),
		map[string]string{
			"accepts_incomplete": "true",
		},
		&ProvisioningRequest{
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
	s, _, err := getTestServer()
	assert.Nil(t, err)
	req, err := getProvisionRequest(
		getDisposableInstanceID(),
		map[string]string{
			"accepts_incomplete": "true",
		},
		&ProvisioningRequest{
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
	s, _, err := getTestServer()
	assert.Nil(t, err)
	req, err := getProvisionRequest(
		getDisposableInstanceID(),
		map[string]string{
			"accepts_incomplete": "true",
		},
		&ProvisioningRequest{
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
	s, _, err := getTestServer()
	assert.Nil(t, err)
	req, err := getProvisionRequest(
		getDisposableInstanceID(),
		map[string]string{
			"accepts_incomplete": "true",
		},
		&ProvisioningRequest{
			ServiceID: fake.ServiceID,
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
	s, _, err := getTestServer()
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	existingInstance := &service.Instance{
		InstanceID: instanceID,
		ServiceID:  fake.ServiceID,
		PlanID:     fake.StandardPlanID,
	}
	err = existingInstance.SetProvisioningParameters(
		&fake.ProvisioningParameters{
			SomeParameter: "foo",
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
		&ProvisioningRequest{
			ServiceID: fake.ServiceID,
			PlanID:    fake.StandardPlanID,
			Parameters: &fake.ProvisioningParameters{
				SomeParameter: "bar",
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
	s, _, err := getTestServer()
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	err = s.store.WriteInstance(&service.Instance{
		InstanceID: instanceID,
		ServiceID:  fake.ServiceID,
		PlanID:     fake.StandardPlanID,
		Status:     service.InstanceStateProvisioned,
	})
	assert.Nil(t, err)
	req, err := getProvisionRequest(
		instanceID,
		map[string]string{
			"accepts_incomplete": "true",
		},
		&ProvisioningRequest{
			ServiceID: fake.ServiceID,
			PlanID:    fake.StandardPlanID,
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
	s, _, err := getTestServer()
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	err = s.store.WriteInstance(&service.Instance{
		InstanceID: instanceID,
		ServiceID:  fake.ServiceID,
		PlanID:     fake.StandardPlanID,
		Status:     service.InstanceStateProvisioning,
	})
	assert.Nil(t, err)
	req, err := getProvisionRequest(
		instanceID,
		map[string]string{
			"accepts_incomplete": "true",
		},
		&ProvisioningRequest{
			ServiceID: fake.ServiceID,
			PlanID:    fake.StandardPlanID,
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusAccepted, rr.Code)
	assert.Equal(t, responseProvisioningAccepted, rr.Body.Bytes())
}

func TestKickOffNewAsyncProvisioning(t *testing.T) {
	s, m, err := getTestServer()
	assert.Nil(t, err)
	validationCalled := false
	m.ProvisioningValidationBehavior = func(service.ProvisioningParameters) error {
		validationCalled = true
		return nil
	}
	instanceID := getDisposableInstanceID()
	req, err := getProvisionRequest(
		instanceID,
		map[string]string{
			"accepts_incomplete": "true",
		},
		&ProvisioningRequest{
			ServiceID: fake.ServiceID,
			PlanID:    fake.StandardPlanID,
		},
	)
	assert.Nil(t, err)
	e := s.asyncEngine.(*fakeAsync.Engine)
	assert.Empty(t, e.SubmittedTasks)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusAccepted, rr.Code)
	assert.True(t, validationCalled)
	assert.Equal(t, 1, len(e.SubmittedTasks))
	assert.Equal(t, responseProvisioningAccepted, rr.Body.Bytes())
}

func getProvisionRequest(
	instanceID string,
	queryParams map[string]string,
	pr *ProvisioningRequest,
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
