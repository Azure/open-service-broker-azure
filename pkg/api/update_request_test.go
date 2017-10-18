package api

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testUpdatingRequest     *UpdatingRequest
	testUpdatingRequestJSON []byte
)

func init() {
	serviceID := "test-service-id"
	planID := "test-plan-id"

	testUpdatingRequest = &UpdatingRequest{
		ServiceID:  serviceID,
		PlanID:     planID,
		Parameters: testArbitraryObject,
	}

	testUpdatingRequestJSONStr := fmt.Sprintf(
		`{
			"service_id":"%s",
			"plan_id":"%s",
			"parameters":%s,
			"previous_values":%s
		}`,
		serviceID,
		planID,
		testArbitraryObjectJSON,
		[]byte(fmt.Sprintf(`{"plan_id":""}`)),
	)
	whitespace := regexp.MustCompile(`\s`)
	testUpdatingRequestJSON = []byte(
		whitespace.ReplaceAllString(testUpdatingRequestJSONStr, ""),
	)
}

func TestGetUpdatingRequestFromJSON(t *testing.T) {
	updatingRequest := &UpdatingRequest{
		Parameters: &ArbitraryType{},
	}
	err := GetUpdatingRequestFromJSON(
		testUpdatingRequestJSON,
		updatingRequest,
	)
	assert.Nil(t, err)
	assert.Equal(t, testUpdatingRequest, updatingRequest)
}

func TestUpdatingRequestToJSON(t *testing.T) {
	json, err := testUpdatingRequest.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, testUpdatingRequestJSON, json)
}
