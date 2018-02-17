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
	bindingDetails := &ArbitraryType{
		Foo: "bat",
	}
	bindingDetailsJSONStr := `{"foo":"bat"}`
	secureBindingDetails := &ArbitraryType{
		Foo: "bat",
	}
	encryptedSecureBindingDetails := []byte(`{"foo":"bat"}`)
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
		Details:                    bindingDetails,
		EncryptedSecureDetails:     encryptedSecureBindingDetails,
		SecureDetails:              secureBindingDetails,
		Created:                    created,
	}

	b64EncryptedBindingParameters := base64.StdEncoding.EncodeToString(
		encryptedBindingParameters,
	)
	b64EncryptedSecureBindingDetails := base64.StdEncoding.EncodeToString(
		encryptedSecureBindingDetails,
	)

	testBindingJSONStr := fmt.Sprintf(
		`{
			"bindingId":"%s",
			"instanceId":"%s",
			"serviceId":"%s",
			"bindingParameters":"%s",
			"status":"%s",
			"statusReason":"%s",
			"details":%s,
			"secureDetails":"%s",
			"created":"%s"
		}`,
		bindingID,
		instanceID,
		serviceID,
		b64EncryptedBindingParameters,
		BindingStateBound,
		statusReason,
		bindingDetailsJSONStr,
		b64EncryptedSecureBindingDetails,
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

func TestEncryptSecureBindingDetails(t *testing.T) {
	binding := Binding{
		SecureDetails: testArbitraryObject,
	}
	var err error
	binding, err = binding.encryptSecureDetails(noopCodec)
	assert.Nil(t, err)
	assert.Equal(
		t,
		testArbitraryObjectJSON,
		binding.EncryptedSecureDetails,
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

func TestDecryptSecureBindingDetails(t *testing.T) {
	binding := Binding{
		EncryptedSecureDetails: testArbitraryObjectJSON,
		SecureDetails:          &ArbitraryType{},
	}
	var err error
	binding, err = binding.decryptSecureDetails(noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testArbitraryObject, binding.SecureDetails)
}
