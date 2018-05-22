package api

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	fakeAsync "github.com/Azure/open-service-broker-azure/pkg/async/fake"
	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/Azure/open-service-broker-azure/pkg/services/fake"
	"github.com/stretchr/testify/assert"
)

func TestProvisioningWithAcceptIncompleteNotSet(t *testing.T) {
	s, _, err := getTestServer("", "")
	assert.Nil(t, err)
	req, err := getProvisionRequest(getDisposableInstanceID(), nil, nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
	assert.Equal(t, responseAsyncRequired, rr.Body.Bytes())
}

func TestProvisioningWithAcceptIncompleteNotTrue(t *testing.T) {
	s, _, err := getTestServer("", "")
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
	s, _, err := getTestServer("", "")
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
	s, _, err := getTestServer("", "")
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
	s, _, err := getTestServer("", "")
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
	s, _, err := getTestServer("", "")
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
	s, _, err := getTestServer("", "")
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	existingInstance := service.Instance{
		InstanceID: instanceID,
		ServiceID:  fake.ServiceID,
		PlanID:     fake.StandardPlanID,
		ProvisioningParameters: service.ProvisioningParameters{
			"someParameter": "foo",
		},
	}
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
			Parameters: map[string]interface{}{
				"someParameter": "bar",
			},
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Equal(t, responseConflict, rr.Body.Bytes())
}

func TestProvisioningWithExistingInstanceWithSameAttributesAndFullyProvisioned(
	t *testing.T,
) {
	s, _, err := getTestServer("", "")
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	err = s.store.WriteInstance(service.Instance{
		InstanceID: instanceID,
		ServiceID:  fake.ServiceID,
		PlanID:     fake.StandardPlanID,
		ProvisioningParameters: service.ProvisioningParameters{
			"someParameter": "foo",
		},
		Status: service.InstanceStateProvisioned,
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

func TestProvisioningWithExistingInstanceWithSameAttributesAndNotFullyProvisioned( // nolint: lll
	t *testing.T,
) {
	s, _, err := getTestServer("", "")
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	err = s.store.WriteInstance(service.Instance{
		InstanceID: instanceID,
		ServiceID:  fake.ServiceID,
		PlanID:     fake.StandardPlanID,
		ProvisioningParameters: service.ProvisioningParameters{
			"someParameter": "foo",
		},
		Status: service.InstanceStateProvisioning,
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
			Parameters: map[string]interface{}{
				"someParameter": "foo",
			},
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusAccepted, rr.Code)
	assert.Equal(t, responseProvisioningAccepted, rr.Body.Bytes())
}

func TestValidatingLocationParameterFails(t *testing.T) {
	s, _, err := getTestServer("", "")
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	req, err := getProvisionRequest(
		instanceID,
		map[string]string{
			"accepts_incomplete": "true",
		},
		&ProvisioningRequest{
			ServiceID: fake.ServiceID,
			PlanID:    fake.StandardPlanID,
			Parameters: map[string]interface{}{
				"location": "upsidedown",
			},
		},
	)
	assert.Nil(t, err)
	e := s.asyncEngine.(*fakeAsync.Engine)
	assert.NotNil(t, e)
	assert.Empty(t, e.SubmittedTasks)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	locationError := service.NewValidationError(
		"location",
		`invalid location: "upsidedown"`,
	)
	responseError := generateValidationFailedResponse(locationError)
	assert.Equal(t, responseError, rr.Body.Bytes())
}

// TODO: krancour: Figure out how to test this now
// func TestModuleSpecificValidationFails(t *testing.T) {
// 	s, m, err := getTestServer("", "")
// 	assert.Nil(t, err)
// 	fooError := service.NewValidationError("foo", "bar")
// 	moduleSpecificValidationCalled := false
// 	m.ServiceManager.ProvisioningValidationBehavior =
// 		func(
// 			service.Plan,
// 			service.ProvisioningParameters,
// 			service.SecureProvisioningParameters,
// 		) error {
// 			moduleSpecificValidationCalled = true
// 			return fooError
// 		}
// 	instanceID := getDisposableInstanceID()
// 	req, err := getProvisionRequest(
// 		instanceID,
// 		map[string]string{
// 			"accepts_incomplete": "true",
// 		},
// 		&ProvisioningRequest{
// 			ServiceID: fake.ServiceID,
// 			PlanID:    fake.StandardPlanID,
// 			Parameters: map[string]interface{}{
// 				"location": "eastus",
// 			},
// 		},
// 	)
// 	assert.Nil(t, err)
// 	e := s.asyncEngine.(*fakeAsync.Engine)
// 	assert.NotNil(t, e)
// 	assert.Empty(t, e.SubmittedTasks)
// 	rr := httptest.NewRecorder()
// 	s.router.ServeHTTP(rr, req)
// 	assert.Equal(t, http.StatusBadRequest, rr.Code)
// 	assert.True(t, moduleSpecificValidationCalled)

// 	responseError := generateValidationFailedResponse(fooError)
// 	assert.Equal(t, responseError, rr.Body.Bytes())
// }

func TestKickOffNewAsyncProvisioning(t *testing.T) {
	s, _, err := getTestServer("", "")
	assert.Nil(t, err)
	instanceID := getDisposableInstanceID()
	req, err := getProvisionRequest(
		instanceID,
		map[string]string{
			"accepts_incomplete": "true",
		},
		&ProvisioningRequest{
			ServiceID: fake.ServiceID,
			PlanID:    fake.StandardPlanID,
			Parameters: map[string]interface{}{
				"location": "eastus",
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
	assert.Equal(t, responseProvisioningAccepted, rr.Body.Bytes())
}

func TestGetLocation(t *testing.T) {
	const defaultLocation = "default-location"
	const location = "test-location"
	testCases := []struct {
		name            string
		svc             service.Service
		defaultLocation string
		location        string
		assertion       func(*testing.T, string)
	}{
		{
			name: `service is not a "root" service type`,
			svc: service.NewService(
				&service.ServiceProperties{
					ParentServiceID: "foo",
				},
				nil,
			),
			location: location,
			assertion: func(t *testing.T, loc string) {
				assert.Equal(t, "", loc)
			},
		},
		{
			name:     "location specified with no default location",
			location: location,
			assertion: func(t *testing.T, loc string) {
				assert.Equal(t, location, loc)
			},
		},
		{
			name:            "location specified with default location",
			location:        location,
			defaultLocation: defaultLocation,
			assertion: func(t *testing.T, loc string) {
				assert.Equal(t, location, loc)
			},
		},
		{
			name: "location not specified with no default location",
			assertion: func(t *testing.T, loc string) {
				assert.Equal(t, "", loc)
			},
		},
	}
	for _, testCase := range testCases {
		if testCase.svc == nil {
			testCase.svc = service.NewService(
				&service.ServiceProperties{},
				nil,
			)
		}
		t.Run(testCase.name, func(t *testing.T) {
			s, _, err := getTestServer(
				testCase.defaultLocation,
				"default-rg",
			)
			assert.Nil(t, err)
			spc := s.getLocation(testCase.svc, testCase.location)
			testCase.assertion(t, spc)
		})
	}
}

func TestGetResourceGroup(t *testing.T) {
	const defaultResourceGroup = "default-rg"
	const resourceGroup = "test-rg"
	testCases := []struct {
		name                 string
		svc                  service.Service
		defaultResourceGroup string
		resourceGroup        string
		assertion            func(*testing.T, string)
	}{
		{
			name: `service is not a "root" service type`,
			svc: service.NewService(
				&service.ServiceProperties{
					ParentServiceID: "foo",
				},
				nil,
			),
			resourceGroup: resourceGroup,
			assertion: func(t *testing.T, rg string) {
				assert.Equal(t, "", rg)
			},
		},
		{
			name:          "resource group specified with no default resource group",
			resourceGroup: resourceGroup,
			assertion: func(t *testing.T, rg string) {
				assert.Equal(t, resourceGroup, rg)
			},
		},
		{
			name:                 "resource group specified with default resource group", // nolint: lll
			resourceGroup:        resourceGroup,
			defaultResourceGroup: defaultResourceGroup,
			assertion: func(t *testing.T, rg string) {
				assert.Equal(t, resourceGroup, rg)
			},
		},
		{
			name: "resource group not specified with no default resource group",
			assertion: func(t *testing.T, rg string) {
				assert.Regexp(
					t,
					regexp.MustCompile(
						`^[\da-f]{8}-[\da-f]{4}-4[\da-f]{3}-[89ab][\da-f]{3}-[\da-f]{12}$`,
					),
					rg,
				)
			},
		},
		{
			name:                 "resource group not specified with default resource group", // nolint: lll
			defaultResourceGroup: defaultResourceGroup,
			assertion: func(t *testing.T, rg string) {
				assert.Equal(t, defaultResourceGroup, rg)
			},
		},
	}
	for _, testCase := range testCases {
		if testCase.svc == nil {
			testCase.svc = service.NewService(
				&service.ServiceProperties{},
				nil,
			)
		}
		t.Run(testCase.name, func(t *testing.T) {
			s, _, err := getTestServer(
				"default-location",
				testCase.defaultResourceGroup,
			)
			assert.Nil(t, err)
			spc := s.getResourceGroup(testCase.svc, testCase.resourceGroup)
			testCase.assertion(t, spc)
		})
	}
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
