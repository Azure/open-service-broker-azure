package service

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testInstance     *Instance
	testInstanceJSON string
)

func init() {
	instanceID := "test-instance-id"
	serviceID := "test-service-id"
	planID := "test-plan-id"
	encodedProvisiongingParameters := `{"foo":"bar"}`
	statusReason := "in-progress"
	encodedProvisiongingContext := `{"baz":"bat"}`

	testInstance = &Instance{
		InstanceID: instanceID,
		ServiceID:  serviceID,
		PlanID:     planID,
		EncodedProvisioningParameters: encodedProvisiongingParameters,
		Status:                     InstanceStateProvisioning,
		StatusReason:               statusReason,
		EncodedProvisioningContext: encodedProvisiongingContext,
	}

	testInstanceJSON = fmt.Sprintf(
		`{
			"instanceId":"%s",
			"serviceId":"%s",
			"planId":"%s",
			"provisioningParameters":%s,
			"status":"%s",
			"statusReason":"%s",
			"provisioningContext":%s
		}`,
		instanceID,
		serviceID,
		planID,
		strconv.Quote(encodedProvisiongingParameters),
		InstanceStateProvisioning,
		statusReason,
		strconv.Quote(encodedProvisiongingContext),
	)
	testInstanceJSON = strings.Replace(testInstanceJSON, " ", "", -1)
	testInstanceJSON = strings.Replace(testInstanceJSON, "\n", "", -1)
	testInstanceJSON = strings.Replace(testInstanceJSON, "\t", "", -1)
}

func TestNewInstanceFromJSONString(t *testing.T) {
	instance, err := NewInstanceFromJSONString(testInstanceJSON)
	assert.Nil(t, err)
	assert.Equal(t, testInstance, instance)
}

func TestInstanceToJSON(t *testing.T) {
	jsonStr, err := testInstance.ToJSONString()
	assert.Nil(t, err)
	assert.Equal(t, testInstanceJSON, jsonStr)
}

func TestSetProvisioningParametersOnInstance(t *testing.T) {
	err := testInstance.SetProvisioningParameters(testArbitraryObject)
	assert.Nil(t, err)
	assert.Equal(
		t,
		testArbitraryObjectJSON,
		testInstance.EncodedProvisioningParameters,
	)
}

func TestGetProvisioningParametersOnInstance(t *testing.T) {
	testInstance.EncodedProvisioningParameters = testArbitraryObjectJSON
	pp := &ArbitraryType{}
	err := testInstance.GetProvisioningParameters(pp)
	assert.Nil(t, err)
	assert.Equal(t, testArbitraryObject, pp)
}

func TestSetProvisioningContextOnInstance(t *testing.T) {
	err := testInstance.SetProvisioningContext(testArbitraryObject)
	assert.Nil(t, err)
	assert.Equal(
		t,
		testInstance.EncodedProvisioningContext,
		testArbitraryObjectJSON,
	)
}

func TestGetProvisioningContextOnInstance(t *testing.T) {
	testInstance.EncodedProvisioningContext = testArbitraryObjectJSON
	pc := &ArbitraryType{}
	err := testInstance.GetProvisioningContext(pc)
	assert.Nil(t, err)
	assert.Equal(t, testArbitraryObject, pc)
}
