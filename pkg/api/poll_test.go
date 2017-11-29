package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azure/azure-service-broker/pkg/service"
	"github.com/stretchr/testify/assert"
)

func TestPollingWithMissingOperation(t *testing.T) {
	s, _, err := getTestServer("", "")
	assert.Nil(t, err)
	req, err := getPollingRequest(getDisposableInstanceID(), "")
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, responseOperationRequired, rr.Body.Bytes())
}

func TestPollingWithInvalidOperation(t *testing.T) {
	s, _, err := getTestServer("", "")
	assert.Nil(t, err)
	req, err := getPollingRequest(getDisposableInstanceID(), "bogus")
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, responseOperationInvalid, rr.Body.Bytes())
}

func TestPollingWithInstanceProvisioning(t *testing.T) {
	s, _, err := getTestServer("", "")
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	err = s.store.WriteInstance(&service.Instance{
		InstanceID: instanceID,
		Status:     service.InstanceStateProvisioning,
	})
	assert.Nil(t, err)
	req, err := getPollingRequest(instanceID, OperationProvisioning)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, responseInProgress, rr.Body.Bytes())
}

func TestPollingWithInstanceProvisioned(t *testing.T) {
	s, _, err := getTestServer("", "")
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	err = s.store.WriteInstance(&service.Instance{
		InstanceID: instanceID,
		Status:     service.InstanceStateProvisioned,
	})
	assert.Nil(t, err)
	req, err := getPollingRequest(instanceID, OperationProvisioning)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, responseSucceeded, rr.Body.Bytes())
}

func TestPollingWithInstanceProvisioningFailed(t *testing.T) {
	s, _, err := getTestServer("", "")
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	err = s.store.WriteInstance(&service.Instance{
		InstanceID: instanceID,
		Status:     service.InstanceStateProvisioningFailed,
	})
	assert.Nil(t, err)
	req, err := getPollingRequest(instanceID, OperationProvisioning)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, responseFailed, rr.Body.Bytes())
}

func TestPollingWithInstanceDeprovisioning(t *testing.T) {
	s, _, err := getTestServer("", "")
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	err = s.store.WriteInstance(&service.Instance{
		InstanceID: instanceID,
		Status:     service.InstanceStateDeprovisioning,
	})
	assert.Nil(t, err)
	req, err := getPollingRequest(instanceID, OperationDeprovisioning)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, responseInProgress, rr.Body.Bytes())
}

func TestPollingWithInstanceGone(t *testing.T) {
	s, _, err := getTestServer("", "")
	assert.Nil(t, err)
	assert.Nil(t, err)
	req, err := getPollingRequest(
		getDisposableInstanceID(),
		OperationDeprovisioning,
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusGone, rr.Code)
	assert.Equal(t, responseEmptyJSON, rr.Body.Bytes())
}

func TestPollingWithInstanceDeprovisioningFailed(t *testing.T) {
	s, _, err := getTestServer("", "")
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	err = s.store.WriteInstance(&service.Instance{
		InstanceID: instanceID,
		Status:     service.InstanceStateDeprovisioningFailed,
	})
	assert.Nil(t, err)
	req, err := getPollingRequest(instanceID, OperationDeprovisioning)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, responseFailed, rr.Body.Bytes())
}

func getPollingRequest(instanceID, operation string) (*http.Request, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("/v2/service_instances/%s/last_operation", instanceID),
		nil,
	)
	if err != nil {
		return nil, err
	}
	if operation != "" {
		q := req.URL.Query()
		q.Add("operation", operation)
		req.URL.RawQuery = q.Encode()
	}
	return req, nil
}
