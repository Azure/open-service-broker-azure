package service

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testBindingRequest     *BindingRequest
	testBindingRequestJSON string
)

func init() {
	serviceID := "test-service-id"
	planID := "test-plan-id"

	testBindingRequest = &BindingRequest{
		ServiceID:  serviceID,
		PlanID:     planID,
		Parameters: testArbitraryObject,
	}

	testBindingRequestJSON = fmt.Sprintf(
		`{
			"service_id":"%s",
			"plan_id":"%s",
			"parameters":%s
		}`,
		serviceID,
		planID,
		testArbitraryObjectJSON,
	)
	whitespace := regexp.MustCompile("\\s")
	testBindingRequestJSON = whitespace.ReplaceAllString(testBindingRequestJSON, "")
}

func TestGetBindingRequestFromJSONString(t *testing.T) {
	bindingRequest := &BindingRequest{
		Parameters: &ArbitraryType{},
	}
	err := GetBindingRequestFromJSONString(
		testBindingRequestJSON,
		bindingRequest,
	)
	assert.Nil(t, err)
	assert.Equal(t, testBindingRequest, bindingRequest)
}

func TestBindingRequestToJSON(t *testing.T) {
	jsonStr, err := testBindingRequest.ToJSONString()
	assert.Nil(t, err)
	assert.Equal(t, testBindingRequestJSON, jsonStr)
}
