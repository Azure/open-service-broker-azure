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
		Parameters: testArbitraryMap,
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

func TestNewBindingRequestFromJSON(t *testing.T) {
	bindingRequest, err := NewBindingRequestFromJSON(
		testBindingRequestJSON,
	)
	assert.Nil(t, err)
	assert.Equal(t, testBindingRequest, bindingRequest)
}

func TestBindingRequestToJSON(t *testing.T) {
	json, err := testBindingRequest.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, testBindingRequestJSON, json)
}
