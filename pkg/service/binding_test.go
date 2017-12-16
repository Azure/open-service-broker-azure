package service

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testBinding     Binding
	testBindingJSON []byte
)

func init() {
	bindingID := "test-binding-id"
	instanceID := "test-instance-id"
	serviceID := "test-service-id"
	encryptedBindingParameters := []byte(`{"foo":"bar"}`)
	bindingParameters := &ArbitraryType{
		Foo: "bar",
	}
	statusReason := "in-progress"
	encryptedBindingContext := []byte(`{"foo":"bat"}`)
	bindingContext := &ArbitraryType{
		Foo: "bat",
	}
	encryptedCredentials := []byte(`{"foo":"baz"}`)
	credentials := &ArbitraryType{
		Foo: "baz",
	}
	created, err := time.Parse(time.RFC3339, "2016-07-22T10:11:55-04:00")
	if err != nil {
		panic(err)
	}

	testBinding = Binding{
		BindingID:                  bindingID,
		InstanceID:                 instanceID,
		ServiceID:                  serviceID,
		EncryptedBindingParameters: encryptedBindingParameters,
		BindingParameters:          bindingParameters,
		Status:                     BindingStateBound,
		StatusReason:               statusReason,
		EncryptedBindingContext:    encryptedBindingContext,
		BindingContext:             bindingContext,
		EncryptedCredentials:       encryptedCredentials,
		Credentials:                credentials,
		Created:                    created,
	}

	b64EncryptedBindingParameters := base64.StdEncoding.EncodeToString(
		encryptedBindingParameters,
	)
	b64EncryptedBindingContext := base64.StdEncoding.EncodeToString(
		encryptedBindingContext,
	)
	b64EncryptedCredentials := base64.StdEncoding.EncodeToString(
		encryptedCredentials,
	)

	testBindingJSONStr := fmt.Sprintf(
		`{
			"bindingId":"%s",
			"instanceId":"%s",
			"serviceId":"%s",
			"bindingParameters":"%s",
			"status":"%s",
			"statusReason":"%s",
			"bindingContext":"%s",
			"credentials":"%s",
			"created":"%s"
		}`,
		bindingID,
		instanceID,
		serviceID,
		b64EncryptedBindingParameters,
		BindingStateBound,
		statusReason,
		b64EncryptedBindingContext,
		b64EncryptedCredentials,
		created.Format(time.RFC3339),
	)
	testBindingJSONStr = strings.Replace(testBindingJSONStr, " ", "", -1)
	testBindingJSONStr = strings.Replace(testBindingJSONStr, "\n", "", -1)
	testBindingJSONStr = strings.Replace(testBindingJSONStr, "\t", "", -1)
	testBindingJSON = []byte(testBindingJSONStr)
}

func TestNewBindingFromJSON(t *testing.T) {
	binding, err := NewBindingFromJSON(
		testBindingJSON,
		&ArbitraryType{},
		&ArbitraryType{},
		&ArbitraryType{},
		noopCodec,
	)
	assert.Nil(t, err)
	assert.Equal(t, testBinding, binding)
}

func TestBindingToJSON(t *testing.T) {
	json, err := testBinding.ToJSON(noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testBindingJSON, json)
}

func TestEncryptBindingParameters(t *testing.T) {
	binding := Binding{
		BindingParameters: testArbitraryObject,
	}
	var err error
	binding, err = binding.encryptBindingParameters(noopCodec)
	assert.Nil(t, err)
	assert.Equal(
		t,
		testArbitraryObjectJSON,
		binding.EncryptedBindingParameters,
	)
}

func TestEncryptBindingContext(t *testing.T) {
	binding := Binding{
		BindingContext: testArbitraryObject,
	}
	var err error
	binding, err = binding.encryptBindingContext(noopCodec)
	assert.Nil(t, err)
	assert.Equal(
		t,
		testArbitraryObjectJSON,
		binding.EncryptedBindingContext,
	)
}

func TestEncryptCredentials(t *testing.T) {
	binding := Binding{
		Credentials: testArbitraryObject,
	}
	var err error
	binding, err = binding.encryptCredentials(noopCodec)
	assert.Nil(t, err)
	assert.Equal(
		t,
		testArbitraryObjectJSON,
		binding.EncryptedCredentials,
	)
}

func TestDecryptBindingParameters(t *testing.T) {
	binding := Binding{
		EncryptedBindingParameters: testArbitraryObjectJSON,
		BindingParameters:          &ArbitraryType{},
	}
	var err error
	binding, err = binding.decryptBindingParameters(noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testArbitraryObject, binding.BindingParameters)
}

func TestDecryptBindingContext(t *testing.T) {
	binding := Binding{
		EncryptedBindingContext: testArbitraryObjectJSON,
		BindingContext:          &ArbitraryType{},
	}
	var err error
	binding, err = binding.decryptBindingContext(noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testArbitraryObject, binding.BindingContext)
}

func TestDecryptCredentials(t *testing.T) {
	binding := Binding{
		EncryptedCredentials: testArbitraryObjectJSON,
		Credentials:          &ArbitraryType{},
	}
	var err error
	binding, err = binding.decryptCredentials(noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testArbitraryObject, binding.Credentials)
}
