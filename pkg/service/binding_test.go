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
	encryptedBindingDetails := []byte(`{"foo":"bat"}`)
	bindingDetails := &ArbitraryType{
		Foo: "bat",
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
		EncryptedDetails:           encryptedBindingDetails,
		Details:                    bindingDetails,
		Created:                    created,
	}

	b64EncryptedBindingParameters := base64.StdEncoding.EncodeToString(
		encryptedBindingParameters,
	)
	b64EncryptedBindingDetails := base64.StdEncoding.EncodeToString(
		encryptedBindingDetails,
	)

	testBindingJSONStr := fmt.Sprintf(
		`{
			"bindingId":"%s",
			"instanceId":"%s",
			"serviceId":"%s",
			"bindingParameters":"%s",
			"status":"%s",
			"statusReason":"%s",
			"details":"%s",
			"created":"%s"
		}`,
		bindingID,
		instanceID,
		serviceID,
		b64EncryptedBindingParameters,
		BindingStateBound,
		statusReason,
		b64EncryptedBindingDetails,
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

func TestEncryptBindingDetails(t *testing.T) {
	binding := Binding{
		Details: testArbitraryObject,
	}
	var err error
	binding, err = binding.encryptDetails(noopCodec)
	assert.Nil(t, err)
	assert.Equal(
		t,
		testArbitraryObjectJSON,
		binding.EncryptedDetails,
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

func TestDecryptBindingDetails(t *testing.T) {
	binding := Binding{
		EncryptedDetails: testArbitraryObjectJSON,
		Details:          &ArbitraryType{},
	}
	var err error
	binding, err = binding.decryptDetails(noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testArbitraryObject, binding.Details)
}
