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
	encryptedProvisiongingParameters := `{"foo":"bar"}`
	statusReason := "in-progress"
	encryptedProvisiongingContext := `{"baz":"bat"}`

	testInstance = &Instance{
		InstanceID: instanceID,
		ServiceID:  serviceID,
		PlanID:     planID,
		EncryptedProvisioningParameters: encryptedProvisiongingParameters,
		Status:                       InstanceStateProvisioning,
		StatusReason:                 statusReason,
		EncryptedProvisioningContext: encryptedProvisiongingContext,
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
		strconv.Quote(encryptedProvisiongingParameters),
		InstanceStateProvisioning,
		statusReason,
		strconv.Quote(encryptedProvisiongingContext),
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
	err := testInstance.SetProvisioningParameters(testArbitraryObject, noopCodec)
	assert.Nil(t, err)
	assert.Equal(
		t,
		testArbitraryObjectJSON,
		testInstance.EncryptedProvisioningParameters,
	)
}

func TestGetProvisioningParametersOnInstance(t *testing.T) {
	testInstance.EncryptedProvisioningParameters = testArbitraryObjectJSON
	pp := &ArbitraryType{}
	err := testInstance.GetProvisioningParameters(pp, noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testArbitraryObject, pp)
}

func TestSetProvisioningContextOnInstance(t *testing.T) {
	err := testInstance.SetProvisioningContext(testArbitraryObject, noopCodec)
	assert.Nil(t, err)
	assert.Equal(
		t,
		testInstance.EncryptedProvisioningContext,
		testArbitraryObjectJSON,
	)
}

func TestGetProvisioningContextOnInstance(t *testing.T) {
	testInstance.EncryptedProvisioningContext = testArbitraryObjectJSON
	pc := &ArbitraryType{}
	err := testInstance.GetProvisioningContext(pc, noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testArbitraryObject, pc)
}
