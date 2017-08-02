package api

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testBindingRequest     *BindingRequest
	testBindingRequestJSON []byte
)

func init() {
	serviceID := "test-service-id"
	planID := "test-plan-id"

	testBindingRequest = &BindingRequest{
		ServiceID:  serviceID,
		PlanID:     planID,
		Parameters: testArbitraryObject,
	}

	testBindingRequestJSONStr := fmt.Sprintf(
		`{
			"service_id":"%s",
			"plan_id":"%s",
			"parameters":%s
		}`,
		serviceID,
		planID,
		testArbitraryObjectJSON,
	)
	whitespace := regexp.MustCompile(`\s`)
	testBindingRequestJSONStr = whitespace.ReplaceAllString(
		testBindingRequestJSONStr,
		"",
	)
	testBindingRequestJSON = []byte(testBindingRequestJSONStr)
}

func TestGetBindingRequestFromJSON(t *testing.T) {
	bindingRequest := &BindingRequest{
		Parameters: &ArbitraryType{},
	}
	err := GetBindingRequestFromJSON(testBindingRequestJSON, bindingRequest)
	assert.Nil(t, err)
	assert.Equal(t, testBindingRequest, bindingRequest)
}

func TestBindingRequestToJSON(t *testing.T) {
	json, err := testBindingRequest.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, testBindingRequestJSON, json)
}
