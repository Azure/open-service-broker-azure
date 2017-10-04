package service

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testInstance     *Instance
	testInstanceJSON []byte
)

func init() {
	instanceID := "test-instance-id"
	serviceID := "test-service-id"
	planID := "test-plan-id"
	encryptedProvisiongingParameters := []byte(`{"foo":"bar"}`)
	statusReason := "in-progress"
	encryptedProvisiongingContext := []byte(`{"baz":"bat"}`)
	created, err := time.Parse(time.RFC3339, "2016-07-22T10:11:55-04:00")
	if err != nil {
		panic(err)
	}

	testInstance = &Instance{
		InstanceID: instanceID,
		ServiceID:  serviceID,
		PlanID:     planID,
		EncryptedProvisioningParameters: encryptedProvisiongingParameters,
		Status:                       InstanceStateProvisioning,
		StatusReason:                 statusReason,
		EncryptedProvisioningContext: encryptedProvisiongingContext,
		Created: created,
	}

	b64EncryptedProvisioningParameters := base64.StdEncoding.EncodeToString(
		encryptedProvisiongingParameters,
	)
	b64EncryptedProvisioningContext := base64.StdEncoding.EncodeToString(
		encryptedProvisiongingContext,
	)

	testInstanceJSONStr := fmt.Sprintf(
		`{
			"instanceId":"%s",
			"serviceId":"%s",
			"planId":"%s",
			"provisioningParameters":"%s",
			"status":"%s",
			"statusReason":"%s",
			"provisioningContext":"%s",
			"created":"%s"
		}`,
		instanceID,
		serviceID,
		planID,
		b64EncryptedProvisioningParameters,
		InstanceStateProvisioning,
		statusReason,
		b64EncryptedProvisioningContext,
		created.Format(time.RFC3339),
	)
	testInstanceJSONStr = strings.Replace(testInstanceJSONStr, " ", "", -1)
	testInstanceJSONStr = strings.Replace(testInstanceJSONStr, "\n", "", -1)
	testInstanceJSONStr = strings.Replace(testInstanceJSONStr, "\t", "", -1)
	testInstanceJSON = []byte(testInstanceJSONStr)
}

func TestNewInstanceFromJSON(t *testing.T) {
	instance, err := NewInstanceFromJSON(testInstanceJSON)
	assert.Nil(t, err)
	assert.Equal(t, testInstance, instance)
}

func TestInstanceToJSON(t *testing.T) {
	json, err := testInstance.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, testInstanceJSON, json)
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
