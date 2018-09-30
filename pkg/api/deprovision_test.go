package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/fake"
	fakeAsync "github.com/deis/async/fake"
	"github.com/stretchr/testify/assert"
)

func TestDeprovisioningWithAcceptIncompleteNotSet(t *testing.T) {
	s, _, err := getTestServer()
	assert.Nil(t, err)
	req, err := getDeprovisionRequest(getDisposableInstanceID(), nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	assert.Equal(t, responseAsyncRequired, rr.Body.Bytes())
}

func TestDeprovisioningWithAcceptIncompleteNotTrue(t *testing.T) {
	s, _, err := getTestServer()
	assert.Nil(t, err)
	req, err := getDeprovisionRequest(
		getDisposableInstanceID(),
		map[string]string{
			"accepts_incomplete": "false",
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	assert.Equal(t, responseAsyncRequired, rr.Body.Bytes())
}

func TestDeprovisioningInstanceThatIsNotFound(t *testing.T) {
	s, _, err := getTestServer()
	assert.Nil(t, err)
	req, err := getDeprovisionRequest(
		getDisposableInstanceID(),
		map[string]string{
			"accepts_incomplete": "true",
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusGone, rr.Code)
	assert.Equal(t, responseEmptyJSON, rr.Body.Bytes())
}

func TestDeprovisioningInstanceThatIsAlreadyDeprovisioning(t *testing.T) {
	s, _, err := getTestServer()
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	err = s.store.WriteInstance(service.Instance{
		InstanceID: instanceID,
		ServiceID:  fake.ServiceID,
		PlanID:     fake.StandardPlanID,
		Status:     service.InstanceStateDeprovisioning,
		Details:    fake.GetEmptyInstanceDetails(),
	})
	assert.Nil(t, err)
	req, err := getDeprovisionRequest(
		instanceID,
		map[string]string{
			"accepts_incomplete": "true",
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusAccepted, rr.Code)
	assert.Equal(t, responseDeprovisioningAccepted, rr.Body.Bytes())
}

func TestDeprovisioningInstanceThatIsStillProvisioning(t *testing.T) {
	s, _, err := getTestServer()
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	err = s.store.WriteInstance(service.Instance{
		InstanceID: instanceID,
		ServiceID:  fake.ServiceID,
		PlanID:     fake.StandardPlanID,
		Status:     service.InstanceStateProvisioning,
		Details:    fake.GetEmptyInstanceDetails(),
	})
	assert.Nil(t, err)
	req, err := getDeprovisionRequest(
		instanceID,
		map[string]string{
			"accepts_incomplete": "true",
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Equal(t, responseEmptyJSON, rr.Body.Bytes())
}

func TestKickOffNewAsyncDeprovisioning(t *testing.T) {
	s, _, err := getTestServer()
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	err = s.store.WriteInstance(service.Instance{
		InstanceID: instanceID,
		ServiceID:  fake.ServiceID,
		PlanID:     fake.StandardPlanID,
		Status:     service.InstanceStateProvisioned,
		Details:    fake.GetEmptyInstanceDetails(),
	})
	assert.Nil(t, err)
	req, err := getDeprovisionRequest(
		instanceID,
		map[string]string{
			"accepts_incomplete": "true",
		},
	)
	assert.Nil(t, err)
	e := s.asyncEngine.(*fakeAsync.Engine)
	assert.NotNil(t, e)
	assert.Empty(t, e.SubmittedTasks)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusAccepted, rr.Code)
	assert.Equal(t, responseDeprovisioningAccepted, rr.Body.Bytes())
	assert.Equal(t, 1, len(e.SubmittedTasks))
}

func getDeprovisionRequest(
	instanceID string,
	queryParams map[string]string,
) (*http.Request, error) {
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("/v2/service_instances/%s", instanceID),
		nil,
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
