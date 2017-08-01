package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthEndpoint(t *testing.T) {
	s, _, err := getTestServer()
	assert.Nil(t, err)
	req, err := getHealthRequest()
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func getHealthRequest() (*http.Request, error) {
	return http.NewRequest(http.MethodGet, "/healthz", nil)
}
