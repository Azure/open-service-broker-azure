package service

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testProvisioningRequest     *ProvisioningRequest
	testProvisioningRequestJSON string
)

func init() {
	serviceID := "test-service-id"
	planID := "test-plan-id"

	testProvisioningRequest = &ProvisioningRequest{
		ServiceID:  serviceID,
		PlanID:     planID,
		Parameters: testArbitraryObject,
	}

	testProvisioningRequestJSON = fmt.Sprintf(
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
	testProvisioningRequestJSON = whitespace.ReplaceAllString(testProvisioningRequestJSON, "")
}

func TestGetProvisioningRequestFromJSONString(t *testing.T) {
	provisioningRequest := &ProvisioningRequest{
		Parameters: &ArbitraryType{},
	}
	err := GetProvisioningRequestFromJSONString(
		testProvisioningRequestJSON,
		provisioningRequest,
	)
	assert.Nil(t, err)
	assert.Equal(t, testProvisioningRequest, provisioningRequest)
}

func TestProvisioningRequestToJSON(t *testing.T) {
	jsonStr, err := testProvisioningRequest.ToJSONString()
	assert.Nil(t, err)
	assert.Equal(t, testProvisioningRequestJSON, jsonStr)
}
