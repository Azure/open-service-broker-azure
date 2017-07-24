package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeprovisioningRejectsIfAcceptIncompleteNotSet(t *testing.T) {
	s, err := getTestServer()
	assert.Nil(t, err)
	req, err := getDeprovisionRequest(uuid.NewV4().String(), nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.deprovision(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestDeprovisioningRejectsIfAcceptIncompleteNotTrue(t *testing.T) {
	s, err := getTestServer()
	assert.Nil(t, err)
	req, err := getDeprovisionRequest(
		uuid.NewV4().String(),
		map[string]string{
			"accepts_incomplete": "false",
		},
	)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.deprovision(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)
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
