package filters

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIVersionFilterWithHeaderMissing(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	handlerCalled := false
	avf := NewAPIVersionFilter()
	avf.GetHandler(func(http.ResponseWriter, *http.Request) {
		handlerCalled = true
	})(rr, req)
	assert.Equal(t, http.StatusPreconditionFailed, rr.Code)
	assert.False(t, handlerCalled)
}

func TestAPIVersionFilterWithHeaderValueIncorrect(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Broker-API-Version", "1.5")
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	handlerCalled := false
	avf := NewAPIVersionFilter()
	avf.GetHandler(func(http.ResponseWriter, *http.Request) {
		handlerCalled = true
	})(rr, req)
	assert.Equal(t, http.StatusPreconditionFailed, rr.Code)
	assert.False(t, handlerCalled)
}

func TestAPIVersionFilterWithHeaderValueValid(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("X-Broker-API-Version", "2.13")
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	handlerCalled := false
	avf := NewAPIVersionFilter()
	avf.GetHandler(func(http.ResponseWriter, *http.Request) {
		handlerCalled = true
	})(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.True(t, handlerCalled)
}
