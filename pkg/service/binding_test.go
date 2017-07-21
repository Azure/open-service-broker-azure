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
	encodedBindingParameters := `{"foo":"bar"}`
	encodedBindingContext := `{"baz":"bat"}`

	testBinding = &Binding{
		BindingID:                bindingID,
		InstanceID:               instanceID,
		EncodedBindingParameters: encodedBindingParameters,
		Status:                BindingStateBinding,
		EncodedBindingContext: encodedBindingContext,
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
		strconv.Quote(encodedBindingParameters),
		BindingStateBinding,
		strconv.Quote(encodedBindingContext),
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
	err := b.SetBindingParameters(testArbitraryObject)
	assert.Nil(t, err)
	assert.Equal(t, testArbitraryObjectJSON, b.EncodedBindingParameters)
}

func TestGetBindingParametersOnBinding(t *testing.T) {
	b := Binding{
		EncodedBindingParameters: testArbitraryObjectJSON,
	}
	bp := &ArbitraryType{}
	err := b.GetBindingParameters(bp)
	assert.Nil(t, err)
	assert.Equal(t, testArbitraryObject, bp)
}

func TestSetBindingContextOnBinding(t *testing.T) {
	b := Binding{}
	err := b.SetBindingContext(testArbitraryObject)
	assert.Nil(t, err)
	assert.Equal(
		t,
		testArbitraryObjectJSON,
		b.EncodedBindingContext,
	)
}

func TestGetBindingContextOnBinding(t *testing.T) {
	b := Binding{
		EncodedBindingContext: testArbitraryObjectJSON,
	}
	bc := &ArbitraryType{}
	err := b.GetBindingContext(bc)
	assert.Nil(t, err)
	assert.Equal(t, testArbitraryObject, bc)
}
