package service

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testBindingResponse     *BindingResponse
	testBindingResponseJSON string
)

func init() {
	testBindingResponse = &BindingResponse{
		Credentials: testArbitraryObject,
	}

	testBindingResponseJSON = fmt.Sprintf(
		`{
			"credentials":%s
		}`,
		testArbitraryObjectJSON,
	)
	whitespace := regexp.MustCompile("\\s")
	testBindingResponseJSON = whitespace.ReplaceAllString(testBindingResponseJSON, "")
}

func TestGetBindingResponseFromJSONString(t *testing.T) {
	bindingResponse := &BindingResponse{
		Credentials: &ArbitraryType{},
	}
	err := GetBindingResponseFromJSONString(
		testBindingResponseJSON,
		bindingResponse,
	)
	assert.Nil(t, err)
	assert.Equal(t, testBindingResponse, bindingResponse)
}

func TestBindingResponseToJSON(t *testing.T) {
	jsonStr, err := testBindingResponse.ToJSONString()
	assert.Nil(t, err)
	assert.Equal(t, testBindingResponseJSON, jsonStr)
}
