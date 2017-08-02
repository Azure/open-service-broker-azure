package api

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testProvisioningRequest     *ProvisioningRequest
	testProvisioningRequestJSON []byte
)

func init() {
	serviceID := "test-service-id"
	planID := "test-plan-id"

	testProvisioningRequest = &ProvisioningRequest{
		ServiceID:  serviceID,
		PlanID:     planID,
		Parameters: testArbitraryObject,
	}

	testProvisioningRequestJSONStr := fmt.Sprintf(
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
	testProvisioningRequestJSON = []byte(
		whitespace.ReplaceAllString(testProvisioningRequestJSONStr, ""),
	)
}

func TestGetProvisioningRequestFromJSON(t *testing.T) {
	provisioningRequest := &ProvisioningRequest{
		Parameters: &ArbitraryType{},
	}
	err := GetProvisioningRequestFromJSON(
		testProvisioningRequestJSON,
		provisioningRequest,
	)
	assert.Nil(t, err)
	assert.Equal(t, testProvisioningRequest, provisioningRequest)
}

func TestProvisioningRequestToJSON(t *testing.T) {
	json, err := testProvisioningRequest.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, testProvisioningRequestJSON, json)
}
