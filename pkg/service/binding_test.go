package service

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testBinding     *Binding
	testBindingJSON string
)

func init() {
	bindingID := "test-binding-id"
	instanceID := "test-instance-id"
	encryptedBindingParameters := `{"foo":"bar"}`
	encryptedBindingContext := `{"baz":"bat"}`

	testBinding = &Binding{
		BindingID:                  bindingID,
		InstanceID:                 instanceID,
		EncryptedBindingParameters: encryptedBindingParameters,
		Status:                  BindingStateBinding,
		EncryptedBindingContext: encryptedBindingContext,
	}

	testBindingJSON = fmt.Sprintf(
		`{
			"bindingId":"%s",
			"instanceId":"%s",
			"bindingParameters":%s,
			"status":"%s",
			"bindingContext":%s
		}`,
		bindingID,
		instanceID,
		strconv.Quote(encryptedBindingParameters),
		BindingStateBinding,
		strconv.Quote(encryptedBindingContext),
	)
	testBindingJSON = strings.Replace(testBindingJSON, " ", "", -1)
	testBindingJSON = strings.Replace(testBindingJSON, "\n", "", -1)
	testBindingJSON = strings.Replace(testBindingJSON, "\t", "", -1)
}

func TestNewBindingFromJSONString(t *testing.T) {
	binding, err := NewBindingFromJSONString(testBindingJSON)
	assert.Nil(t, err)
	assert.Equal(t, testBinding, binding)
}

func TestBindingToJSON(t *testing.T) {
	jsonStr, err := testBinding.ToJSONString()
	assert.Nil(t, err)
	assert.Equal(t, testBindingJSON, jsonStr)
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
