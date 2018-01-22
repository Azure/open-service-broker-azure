package basic

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Azure/open-service-broker-azure/pkg/api/filters/authenticator"
	"github.com/stretchr/testify/assert"
)

const (
	testUsername = "user"
	testPassword = "password"
)

func TestAuthHeaderMissing(t *testing.T) {
	a := getTestAuthenticator()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.Nil(t, err)
	rr := httptest.NewRecorder()
	handlerCalled := false
	a.Execute(func(http.ResponseWriter, *http.Request) {
		handlerCalled = true
	})(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.False(t, handlerCalled)
}

func TestAuthHeaderNotBasic(t *testing.T) {
	a := getTestAuthenticator()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Digest foo")
	rr := httptest.NewRecorder()
	handlerCalled := false
	a.Execute(func(http.ResponseWriter, *http.Request) {
		handlerCalled = true
	})(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.False(t, handlerCalled)
}

func TestAuthUsernamePasswordNotBase64(t *testing.T) {
	a := getTestAuthenticator()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.Nil(t, err)
	req.Header.Add("Authorization", "Basic foo")
	rr := httptest.NewRecorder()
	handlerCalled := false
	a.Execute(func(http.ResponseWriter, *http.Request) {
		handlerCalled = true
	})(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.False(t, handlerCalled)
}

func TestAuthUsernamePasswordInvalid(t *testing.T) {
	a := getTestAuthenticator()
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
	a.Execute(func(http.ResponseWriter, *http.Request) {
		handlerCalled = true
	})(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.False(t, handlerCalled)
}

func TestAuthUsernamePasswordValid(t *testing.T) {
	a := getTestAuthenticator()
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
	a.Execute(func(http.ResponseWriter, *http.Request) {
		handlerCalled = true
	})(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.True(t, handlerCalled)
}

func getTestAuthenticator() authenticator.Authenticator {
	return NewAuthenticator(testUsername, testPassword)
}
