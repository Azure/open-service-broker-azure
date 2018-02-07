package filters

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azure/open-service-broker-azure/pkg/http/filter"
	"github.com/stretchr/testify/assert"
)

const (
	testUsername = "user"
	testPassword = "password"
)

func TestBasicAuthFilterWithHeaderMissing(t *testing.T) {
	a := getTestBasicAuthFilter()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	handlerCalled := false
	a.GetHandler(func(http.ResponseWriter, *http.Request) {
		handlerCalled = true
	})(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.False(t, handlerCalled)
}

func TestBasicAuthFilterWithHeaderNotBasic(t *testing.T) {
	a := getTestBasicAuthFilter()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Digest foo")
	rr := httptest.NewRecorder()
	handlerCalled := false
	a.GetHandler(func(http.ResponseWriter, *http.Request) {
		handlerCalled = true
	})(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.False(t, handlerCalled)
}

func TestBasicAuthFilterWithUsernamePasswordNotBase64(t *testing.T) {
	a := getTestBasicAuthFilter()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Basic foo")
	rr := httptest.NewRecorder()
	handlerCalled := false
	a.GetHandler(func(http.ResponseWriter, *http.Request) {
		handlerCalled = true
	})(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.False(t, handlerCalled)
}

func TestBasicAuthFilterWithUsernamePasswordInvalid(t *testing.T) {
	a := getTestBasicAuthFilter()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.Nil(t, err)
	usernameAndPassword := fmt.Sprintf("%s:%s", "foo", "bar")
	b64UsernameAndPassword := base64.StdEncoding.EncodeToString(
		[]byte(usernameAndPassword),
	)
	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Basic %s", b64UsernameAndPassword),
	)
	rr := httptest.NewRecorder()
	handlerCalled := false
	a.GetHandler(func(http.ResponseWriter, *http.Request) {
		handlerCalled = true
	})(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.False(t, handlerCalled)
}

func TestBasicAuthFilterWithUsernamePasswordValid(t *testing.T) {
	a := getTestBasicAuthFilter()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.Nil(t, err)
	usernameAndPassword := fmt.Sprintf("%s:%s", testUsername, testPassword)
	b64UsernameAndPassword := base64.StdEncoding.EncodeToString(
		[]byte(usernameAndPassword),
	)
	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Basic %s", b64UsernameAndPassword),
	)
	rr := httptest.NewRecorder()
	handlerCalled := false
	a.GetHandler(func(http.ResponseWriter, *http.Request) {
		handlerCalled = true
	})(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.True(t, handlerCalled)
}

func getTestBasicAuthFilter() filter.Filter {
	return NewBasicAuthFilter(testUsername, testPassword)
}
