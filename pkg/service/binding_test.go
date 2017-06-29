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
	encodedBindingResult := `{"baz":"bat"}`

	testBinding = &Binding{
		BindingID:                bindingID,
		InstanceID:               instanceID,
		EncodedBindingParameters: encodedBindingParameters,
		Status:               BindingStateBinding,
		EncodedBindingResult: encodedBindingResult,
	}

	testBindingJSON = fmt.Sprintf(
		`{
			"bindingId":"%s",
			"instanceId":"%s",
			"bindingParameters":%s,
			"status":"%s",
			"bindingResult":%s
		}`,
		bindingID,
		instanceID,
		strconv.Quote(encodedBindingParameters),
		BindingStateBinding,
		strconv.Quote(encodedBindingResult),
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

func TestSetBindingResultOnBinding(t *testing.T) {
	b := Binding{}
	err := b.SetBindingResult(testArbitraryObject)
	assert.Nil(t, err)
	assert.Equal(
		t,
		testArbitraryObjectJSON,
		b.EncodedBindingResult,
	)
}

func TestGetBindingResultOnBinding(t *testing.T) {
	b := Binding{
		EncodedBindingResult: testArbitraryObjectJSON,
	}
	br := &ArbitraryType{}
	err := b.GetBindingResult(br)
	assert.Nil(t, err)
	assert.Equal(t, testArbitraryObject, br)
}
