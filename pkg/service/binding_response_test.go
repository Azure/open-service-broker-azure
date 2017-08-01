package service

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testBindingResponse     *BindingResponse
	testBindingResponseJSON []byte
)

func init() {
	testBindingResponse = &BindingResponse{
		Credentials: testArbitraryObject,
	}

	testBindingResponseJSONStr := fmt.Sprintf(
		`{
			"credentials":%s
		}`,
		testArbitraryObjectJSON,
	)
	whitespace := regexp.MustCompile(`\s`)
	testBindingResponseJSON = []byte(
		whitespace.ReplaceAllString(testBindingResponseJSONStr, ""),
	)
}

func TestGetBindingResponseFromJSON(t *testing.T) {
	bindingResponse := &BindingResponse{
		Credentials: &ArbitraryType{},
	}
	err := GetBindingResponseFromJSON(testBindingResponseJSON, bindingResponse)
	assert.Nil(t, err)
	assert.Equal(t, testBindingResponse, bindingResponse)
}

func TestBindingResponseToJSON(t *testing.T) {
	json, err := testBindingResponse.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, testBindingResponseJSON, json)
}
