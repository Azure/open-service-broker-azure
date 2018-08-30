package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/fake"
	fakeAsync "github.com/krancour/async/fake"
	"github.com/stretchr/testify/assert"
)

func TestUpdatingWithAcceptIncompleteNotSet(t *testing.T) {
	s, _, err := getTestServer()
	assert.Nil(t, err)
	req, err := getUpdateRequest(getDisposableInstanceID(), nil, nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	assert.Equal(t, responseAsyncRequired, rr.Body.Bytes())
}

func TestUpdatingWithAcceptIncompleteNotTrue(t *testing.T) {
	s, _, err := getTestServer()
	assert.Nil(t, err)
	req, err := getUpdateRequest(
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

func TestUpdatingWithMissingServiceID(t *testing.T) {
	s, _, err := getTestServer()
	assert.Nil(t, err)
	req, err := getUpdateRequest(
		getDisposableInstanceID(),
		map[string]string{
			"accepts_incomplete": "true",
		},
		&UpdatingRequest{
			ServiceID: "",
			PlanID:    fake.StandardPlanID,
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, responseServiceIDRequired, rr.Body.Bytes())
}

func TestUpdatingWithInvalidServiceID(t *testing.T) {
	s, _, err := getTestServer()
	assert.Nil(t, err)
	req, err := getUpdateRequest(
		getDisposableInstanceID(),
		map[string]string{
			"accepts_incomplete": "true",
		},
		&UpdatingRequest{
			ServiceID: getDisposableServiceID(),
			PlanID:    fake.StandardPlanID,
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, responseInvalidServiceID, rr.Body.Bytes())
}

func TestUpdatingWithInvalidPlanID(t *testing.T) {
	s, _, err := getTestServer()
	assert.Nil(t, err)
	req, err := getUpdateRequest(
		getDisposableInstanceID(),
		map[string]string{
			"accepts_incomplete": "true",
		},
		&UpdatingRequest{
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

func TestUpdatingWithExistingInstanceWithSameAttributesAndFullyProvisioned(
	t *testing.T,
) {
	s, _, err := getTestServer()
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	err = s.store.WriteInstance(service.Instance{
		InstanceID: instanceID,
		ServiceID:  fake.ServiceID,
		PlanID:     fake.StandardPlanID,
		ProvisioningParameters: &service.ProvisioningParameters{
			Parameters: service.Parameters{
				Schema: &service.InputParametersSchema{
					PropertySchemas: map[string]service.PropertySchema{
						"someParameter": &service.StringPropertySchema{},
					},
				},
				Data: map[string]interface{}{
					"someParameter": "foo",
				},
			},
		},
		Status: service.InstanceStateProvisioned,
	})
	assert.Nil(t, err)
	req, err := getUpdateRequest(
		instanceID,
		map[string]string{
			"accepts_incomplete": "true",
		},
		&UpdatingRequest{
			ServiceID: fake.ServiceID,
			PlanID:    fake.StandardPlanID,
			Parameters: map[string]interface{}{
				"someParameter": "foo",
			},
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, responseEmptyJSON, rr.Body.Bytes())
}

func TestUpdatingWithExistingInstanceWithSameAttributesAndNotFullyProvisioned( // nolint: lll
	t *testing.T,
) {
	s, _, err := getTestServer()
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	err = s.store.WriteInstance(service.Instance{
		InstanceID: instanceID,
		ServiceID:  fake.ServiceID,
		PlanID:     fake.StandardPlanID,
		ProvisioningParameters: &service.ProvisioningParameters{
			Parameters: service.Parameters{
				Schema: &service.InputParametersSchema{
					PropertySchemas: map[string]service.PropertySchema{
						"someParameter": &service.StringPropertySchema{},
					},
				},
				Data: map[string]interface{}{
					"someParameter": "foo",
				},
			},
		},
		Status: service.InstanceStateProvisioning,
	})
	assert.Nil(t, err)
	req, err := getUpdateRequest(
		instanceID,
		map[string]string{
			"accepts_incomplete": "true",
		},
		&UpdatingRequest{
			ServiceID: fake.ServiceID,
			PlanID:    fake.StandardPlanID,
			Parameters: map[string]interface{}{
				"someParameter": "foo",
			},
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusAccepted, rr.Code)
	assert.Equal(t, responseUpdatingAccepted, rr.Body.Bytes())
}

func TestUpdatingWithExistingInstanceWithSameAttributesAndNotFullyUpdated( // nolint: lll
	t *testing.T,
) {
	s, _, err := getTestServer()
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	err = s.store.WriteInstance(service.Instance{
		InstanceID: instanceID,
		ServiceID:  fake.ServiceID,
		PlanID:     fake.StandardPlanID,
		UpdatingParameters: &service.ProvisioningParameters{
			Parameters: service.Parameters{
				Schema: &service.InputParametersSchema{
					PropertySchemas: map[string]service.PropertySchema{
						"someParameter": &service.StringPropertySchema{},
					},
				},
				Data: map[string]interface{}{
					"someParameter": "foo",
				},
			},
		},
		Status: service.InstanceStateUpdating,
	})
	assert.Nil(t, err)
	req, err := getUpdateRequest(
		instanceID,
		map[string]string{
			"accepts_incomplete": "true",
		},
		&UpdatingRequest{
			ServiceID: fake.ServiceID,
			PlanID:    fake.StandardPlanID,
			Parameters: map[string]interface{}{
				"someParameter": "foo",
			},
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusAccepted, rr.Code)
	assert.Equal(t, responseUpdatingAccepted, rr.Body.Bytes())
}

func TestKickOffNewAsyncUpdating(t *testing.T) {
	s, _, err := getTestServer()
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	err = s.store.WriteInstance(service.Instance{
		InstanceID: instanceID,
		ServiceID:  fake.ServiceID,
		PlanID:     fake.StandardPlanID,
		Status:     service.InstanceStateProvisioned,
	})
	assert.Nil(t, err)
	req, err := getUpdateRequest(
		instanceID,
		map[string]string{
			"accepts_incomplete": "true",
		},
		&UpdatingRequest{
			ServiceID: fake.ServiceID,
			PlanID:    fake.StandardPlanID,
			Parameters: map[string]interface{}{
				"someParameter": "fake",
			},
		},
	)
	assert.Nil(t, err)
	e := s.asyncEngine.(*fakeAsync.Engine)
	assert.NotNil(t, e)
	assert.Empty(t, e.SubmittedTasks)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusAccepted, rr.Code)
	assert.Equal(t, 1, len(e.SubmittedTasks))
	assert.Equal(t, responseUpdatingAccepted, rr.Body.Bytes())
}

func getUpdateRequest(
	instanceID string,
	queryParams map[string]string,
	ur *UpdatingRequest,
) (*http.Request, error) {
	var body []byte
	if ur != nil {
		var err error
		body, err = ur.ToJSON()
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(
		http.MethodPatch,
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
