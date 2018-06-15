package api

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicAuthUsernameNotSpecified(t *testing.T) {
	err := os.Setenv("BASIC_AUTH_USERNAME", "")
	assert.Nil(t, err)
	err = os.Setenv("SECURITY_USER_NAME", "")
	assert.Nil(t, err)
	_, err = GetBasicAuthConfig()
	assert.IsType(t, &errBasicAuthUsernameNotSpecified{}, err)
}

func TestBasicAuthPasswordNotSpecified(t *testing.T) {
	err := os.Setenv("BASIC_AUTH_USERNAME", "foo")
	assert.Nil(t, err)
	err = os.Setenv("BASIC_AUTH_PASSWORD", "")
	assert.Nil(t, err)
	err = os.Setenv("SECURITY_USER_PASSWORD", "")
	assert.Nil(t, err)
	_, err = GetBasicAuthConfig()
	assert.IsType(t, &errBasicAuthPasswordNotSpecified{}, err)
}

func TestBasicAuthUsernamePasswordPrecedence(t *testing.T) {
	err := os.Setenv("BASIC_AUTH_USERNAME", "foo")
	assert.Nil(t, err)
	err = os.Setenv("SECURITY_USER_NAME", "foo2")
	assert.Nil(t, err)
	err = os.Setenv("BASIC_AUTH_PASSWORD", "bar")
	assert.Nil(t, err)
	err = os.Setenv("SECURITY_USER_PASSWORD", "bar2")
	assert.Nil(t, err)
	bac, err := GetBasicAuthConfig()
	assert.Nil(t, err)
	assert.Equal(t, "foo", bac.GetUsername())
	assert.Equal(t, "bar", bac.GetPassword())
}

func TestBasicAuthUsernamePasswordFallback(t *testing.T) {
	err := os.Setenv("BASIC_AUTH_USERNAME", "")
	assert.Nil(t, err)
	err = os.Setenv("SECURITY_USER_NAME", "foo2")
	assert.Nil(t, err)
	err = os.Setenv("BASIC_AUTH_PASSWORD", "")
	assert.Nil(t, err)
	err = os.Setenv("SECURITY_USER_PASSWORD", "bar2")
	assert.Nil(t, err)
	bac, err := GetBasicAuthConfig()
	assert.Nil(t, err)
	assert.Equal(t, "foo2", bac.GetUsername())
	assert.Equal(t, "bar2", bac.GetPassword())
}
