package service

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testBinding     *Binding
	testBindingJSON []byte
)

func init() {
	bindingID := "test-binding-id"
	instanceID := "test-instance-id"
	encryptedBindingParameters := []byte(`{"foo":"bar"}`)
	statusReason := "in-progress"
	encryptedBindingContext := []byte(`{"baz":"bat"}`)
	encryptedCredentials := []byte(`{"password":"12345"}`)

	testBinding = &Binding{
		BindingID:                  bindingID,
		InstanceID:                 instanceID,
		EncryptedBindingParameters: encryptedBindingParameters,
		Status:                  BindingStateBound,
		StatusReason:            statusReason,
		EncryptedBindingContext: encryptedBindingContext,
		EncryptedCredentials:    encryptedCredentials,
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
			"bindingParameters":"%s",
			"status":"%s",
			"statusReason":"%s",
			"bindingContext":"%s",
			"credentials":"%s"
		}`,
		bindingID,
		instanceID,
		b64EncryptedBindingParameters,
		BindingStateBound,
		statusReason,
		b64EncryptedBindingContext,
		b64EncryptedCredentials,
	)
	testBindingJSONStr = strings.Replace(testBindingJSONStr, " ", "", -1)
	testBindingJSONStr = strings.Replace(testBindingJSONStr, "\n", "", -1)
	testBindingJSONStr = strings.Replace(testBindingJSONStr, "\t", "", -1)
	testBindingJSON = []byte(testBindingJSONStr)
}

func TestNewBindingFromJSON(t *testing.T) {
	binding, err := NewBindingFromJSON(testBindingJSON)
	assert.Nil(t, err)
	assert.Equal(t, testBinding, binding)
}

func TestBindingToJSON(t *testing.T) {
	json, err := testBinding.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, testBindingJSON, json)
}

func TestSetBindingParametersOnBinding(t *testing.T) {
	b := Binding{}
	err := b.SetBindingParameters(testArbitraryObject, noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testArbitraryObjectJSON, b.EncryptedBindingParameters)
}

func TestGetBindingParametersOnBinding(t *testing.T) {
	b := Binding{
		EncryptedBindingParameters: testArbitraryObjectJSON,
	}
	bp := &ArbitraryType{}
	err := b.GetBindingParameters(bp, noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testArbitraryObject, bp)
}

func TestSetBindingContextOnBinding(t *testing.T) {
	b := Binding{}
	err := b.SetBindingContext(testArbitraryObject, noopCodec)
	assert.Nil(t, err)
	assert.Equal(
		t,
		testArbitraryObjectJSON,
		b.EncryptedBindingContext,
	)
}

func TestGetBindingContextOnBinding(t *testing.T) {
	b := Binding{
		EncryptedBindingContext: testArbitraryObjectJSON,
	}
	bc := &ArbitraryType{}
	err := b.GetBindingContext(bc, noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testArbitraryObject, bc)
}
