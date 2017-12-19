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
	testInstance     Instance
	testInstanceJSON []byte
)

func init() {
	instanceID := "test-instance-id"
	serviceID := "test-service-id"
	planID := "test-plan-id"
	location := "test-location"
	resourceGroup := "test-rg"
	tagKey := "foo"
	tagVal := "bar"
	provisioningParameters := &ArbitraryType{
		Foo: "bar",
	}
	encryptedProvisiongingParameters := []byte(`{"foo":"bar"}`)
	updatingParameters := &ArbitraryType{
		Foo: "bat",
	}
	encryptedUpdatingParameters := []byte(`{"foo":"bat"}`)
	statusReason := "in-progress"
	details := &ArbitraryType{
		Foo: "baz",
	}
	encryptedProvisiongingContext := []byte(`{"foo":"baz"}`)
	created, err := time.Parse(time.RFC3339, "2016-07-22T10:11:55-04:00")
	if err != nil {
		panic(err)
	}

	testInstance = Instance{
		InstanceID: instanceID,
		ServiceID:  serviceID,
		PlanID:     planID,
		EncryptedProvisioningParameters: encryptedProvisiongingParameters,
		ProvisioningParameters:          provisioningParameters,
		EncryptedUpdatingParameters:     encryptedUpdatingParameters,
		UpdatingParameters:              updatingParameters,
		Status:                          InstanceStateProvisioning,
		StatusReason:                    statusReason,
		Location:                        location,
		ResourceGroup:                   resourceGroup,
		Tags:                            map[string]string{tagKey: tagVal},
		EncryptedDetails:                encryptedProvisiongingContext,
		Details:                         details,
		Created:                         created,
	}

	b64EncryptedProvisioningParameters := base64.StdEncoding.EncodeToString(
		encryptedProvisiongingParameters,
	)
	b64EncryptedUpdatingParameters := base64.StdEncoding.EncodeToString(
		encryptedUpdatingParameters,
	)
	b64EncryptedDetails := base64.StdEncoding.EncodeToString(
		encryptedProvisiongingContext,
	)

	testInstanceJSONStr := fmt.Sprintf(
		`{
			"instanceId":"%s",
			"serviceId":"%s",
			"planId":"%s",
			"provisioningParameters":"%s",
			"updatingParameters":"%s",
			"status":"%s",
			"statusReason":"%s",
			"location":"%s",
			"resourceGroup":"%s",
			"tags":{"%s":"%s"},
			"details":"%s",
			"created":"%s"
		}`,
		instanceID,
		serviceID,
		planID,
		b64EncryptedProvisioningParameters,
		b64EncryptedUpdatingParameters,
		InstanceStateProvisioning,
		statusReason,
		location,
		resourceGroup,
		tagKey,
		tagVal,
		b64EncryptedDetails,
		created.Format(time.RFC3339),
	)
	testInstanceJSONStr = strings.Replace(testInstanceJSONStr, " ", "", -1)
	testInstanceJSONStr = strings.Replace(testInstanceJSONStr, "\n", "", -1)
	testInstanceJSONStr = strings.Replace(testInstanceJSONStr, "\t", "", -1)
	testInstanceJSON = []byte(testInstanceJSONStr)
}

func TestNewInstanceFromJSON(t *testing.T) {
	instance, err := NewInstanceFromJSON(
		testInstanceJSON,
		&ArbitraryType{},
		&ArbitraryType{},
		&ArbitraryType{},
		noopCodec,
	)
	assert.Nil(t, err)
	assert.Equal(t, testInstance, instance)
}

func TestInstanceToJSON(t *testing.T) {
	json, err := testInstance.ToJSON(noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testInstanceJSON, json)
}

func TestEncryptProvisioningParameters(t *testing.T) {
	instance := Instance{
		ProvisioningParameters: testArbitraryObject,
	}
	var err error
	instance, err = instance.encryptProvisioningParameters(noopCodec)
	assert.Nil(t, err)
	assert.Equal(
		t,
		testArbitraryObjectJSON,
		instance.EncryptedProvisioningParameters,
	)
}

func TestEncryptUpdatingParameters(t *testing.T) {
	instance := Instance{
		UpdatingParameters: testArbitraryObject,
	}
	var err error
	instance, err = instance.encryptUpdatingParameters(noopCodec)
	assert.Nil(t, err)
	assert.Equal(
		t,
		testArbitraryObjectJSON,
		instance.EncryptedUpdatingParameters,
	)
}

func TestEncryptDetails(t *testing.T) {
	instance := Instance{
		Details: testArbitraryObject,
	}
	var err error
	instance, err = instance.encryptDetails(noopCodec)
	assert.Nil(t, err)
	assert.Equal(
		t,
		instance.EncryptedDetails,
		testArbitraryObjectJSON,
	)
}

func TestDecryptProvisioningParameters(t *testing.T) {
	instance := Instance{
		EncryptedProvisioningParameters: testArbitraryObjectJSON,
		ProvisioningParameters:          &ArbitraryType{},
	}
	var err error
	instance, err = instance.decryptProvisioningParameters(noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testArbitraryObject, instance.ProvisioningParameters)
}

func TestDecryptUpdatingParameters(t *testing.T) {
	instance := Instance{
		EncryptedUpdatingParameters: testArbitraryObjectJSON,
		UpdatingParameters:          &ArbitraryType{},
	}
	var err error
	instance, err = instance.decryptUpdatingParameters(noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testArbitraryObject, instance.UpdatingParameters)
}

func TestDecryptDetails(t *testing.T) {
	instance := Instance{
		EncryptedDetails: testArbitraryObjectJSON,
		Details:          &ArbitraryType{},
	}
	var err error
	instance, err = instance.decryptDetails(noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testArbitraryObject, instance.Details)
}
