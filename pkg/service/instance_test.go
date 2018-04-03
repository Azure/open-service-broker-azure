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
	testSecureProvisioningParameters = SecureProvisioningParameters{
		"foo": fooValue,
	}
	testSecureInstanceDetails = SecureInstanceDetails{
		"foo": fooValue,
	}
	testInstance     Instance
	testInstanceJSON []byte
)

func init() {
	instanceID := "test-instance-id"
	alias := "test-alias"
	serviceID := "test-service-id"
	planID := "test-plan-id"
	location := "test-location"
	resourceGroup := "test-rg"
	parentAlias := "test-parent-alias"
	tagKey := "foo"
	tagVal := "bar"
	provisioningParameters := ProvisioningParameters{
		"foo": "bar",
	}
	provisioningParametersJSONStr := []byte(`{"foo":"bar"}`)
	secureProvisioningParameters := SecureProvisioningParameters{
		"foo": "bar",
	}
	encryptedSecureProvisiongingParameters := []byte(`{"foo":"bar"}`)
	updatingParameters := ProvisioningParameters{
		"foo": "bat",
	}
	updatingParametersJSONStr := []byte(`{"foo":"bat"}`)
	secureUpdatingParameters := SecureProvisioningParameters{
		"foo": "bat",
	}
	encryptedSecureUpdatingParameters := []byte(`{"foo":"bat"}`)
	statusReason := "in-progress"
	details := InstanceDetails{
		"foo": "baz",
	}
	detailsJSONStr := `{"foo":"baz"}`
	secureDetails := SecureInstanceDetails{
		"foo": "baz",
	}
	encryptedSecureDetails := []byte(`{"foo":"baz"}`)
	created, err := time.Parse(time.RFC3339, "2016-07-22T10:11:55-04:00")
	if err != nil {
		panic(err)
	}

	testInstance = Instance{
		InstanceID:                            instanceID,
		Alias:                                 alias,
		ServiceID:                             serviceID,
		PlanID:                                planID,
		ProvisioningParameters:                provisioningParameters,
		EncryptedSecureProvisioningParameters: encryptedSecureProvisiongingParameters, // nolint: lll
		SecureProvisioningParameters:          secureProvisioningParameters,
		UpdatingParameters:                    updatingParameters,
		EncryptedSecureUpdatingParameters:     encryptedSecureUpdatingParameters,
		SecureUpdatingParameters:              secureUpdatingParameters,
		Status:                 InstanceStateProvisioning,
		StatusReason:           statusReason,
		Location:               location,
		ResourceGroup:          resourceGroup,
		ParentAlias:            parentAlias,
		Tags:                   map[string]string{tagKey: tagVal},
		Details:                details,
		EncryptedSecureDetails: encryptedSecureDetails,
		SecureDetails:          secureDetails,
		Created:                created,
	}

	b64EncryptedSecureProvisioningParameters := base64.StdEncoding.EncodeToString(
		encryptedSecureProvisiongingParameters,
	)
	b64EncryptedSecureUpdatingParameters := base64.StdEncoding.EncodeToString(
		encryptedSecureUpdatingParameters,
	)
	b64EncryptedSecureDetails := base64.StdEncoding.EncodeToString(
		encryptedSecureDetails,
	)

	testInstanceJSONStr := fmt.Sprintf(
		`{
			"instanceId":"%s",
			"alias":"%s",
			"serviceId":"%s",
			"planId":"%s",
			"provisioningParameters":%s,
			"secureProvisioningParameters":"%s",
			"updatingParameters":%s,
			"secureUpdatingParameters":"%s",
			"status":"%s",
			"statusReason":"%s",
			"location":"%s",
			"resourceGroup":"%s",
			"parentAlias":"%s",
			"tags":{"%s":"%s"},
			"details":%s,
			"secureDetails":"%s",
			"created":"%s"
		}`,
		instanceID,
		alias,
		serviceID,
		planID,
		provisioningParametersJSONStr,
		b64EncryptedSecureProvisioningParameters,
		updatingParametersJSONStr,
		b64EncryptedSecureUpdatingParameters,
		InstanceStateProvisioning,
		statusReason,
		location,
		resourceGroup,
		parentAlias,
		tagKey,
		tagVal,
		detailsJSONStr,
		b64EncryptedSecureDetails,
		created.Format(time.RFC3339),
	)
	testInstanceJSONStr = strings.Replace(testInstanceJSONStr, " ", "", -1)
	testInstanceJSONStr = strings.Replace(testInstanceJSONStr, "\n", "", -1)
	testInstanceJSONStr = strings.Replace(testInstanceJSONStr, "\t", "", -1)
	testInstanceJSON = []byte(testInstanceJSONStr)
}

func TestNewInstanceFromJSON(t *testing.T) {
	instance, err := NewInstanceFromJSON(testInstanceJSON, noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testInstance, instance)
}

func TestInstanceToJSON(t *testing.T) {
	json, err := testInstance.ToJSON(noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testInstanceJSON, json)
}

func TestEncryptSecureProvisioningParameters(t *testing.T) {
	instance := Instance{
		SecureProvisioningParameters: testSecureProvisioningParameters,
	}
	var err error
	instance, err = instance.encryptSecureProvisioningParameters(noopCodec)
	assert.Nil(t, err)
	assert.Equal(
		t,
		testArbitraryObjectJSON,
		instance.EncryptedSecureProvisioningParameters,
	)
}

func TestEncryptSecureUpdatingParameters(t *testing.T) {
	instance := Instance{
		SecureUpdatingParameters: testSecureProvisioningParameters,
	}
	var err error
	instance, err = instance.encryptSecureUpdatingParameters(noopCodec)
	assert.Nil(t, err)
	assert.Equal(
		t,
		testArbitraryObjectJSON,
		instance.EncryptedSecureUpdatingParameters,
	)
}

func TestEncryptSecureDetails(t *testing.T) {
	instance := Instance{
		SecureDetails: testSecureInstanceDetails,
	}
	var err error
	instance, err = instance.encryptSecureDetails(noopCodec)
	assert.Nil(t, err)
	assert.Equal(
		t,
		testArbitraryObjectJSON,
		instance.EncryptedSecureDetails,
	)
}

func TestDecryptSecureProvisioningParameters(t *testing.T) {
	instance := Instance{
		EncryptedSecureProvisioningParameters: testArbitraryObjectJSON,
		SecureProvisioningParameters:          SecureProvisioningParameters{},
	}
	var err error
	instance, err = instance.decryptSecureProvisioningParameters(noopCodec)
	assert.Nil(t, err)
	assert.Equal(
		t,
		testSecureProvisioningParameters,
		instance.SecureProvisioningParameters,
	)
}

func TestDecryptUpdatingParameters(t *testing.T) {
	instance := Instance{
		EncryptedSecureUpdatingParameters: testArbitraryObjectJSON,
		SecureUpdatingParameters:          SecureProvisioningParameters{},
	}
	var err error
	instance, err = instance.decryptSecureUpdatingParameters(noopCodec)
	assert.Nil(t, err)
	assert.Equal(
		t,
		testSecureProvisioningParameters,
		instance.SecureUpdatingParameters,
	)
}

func TestDecryptSecureDetails(t *testing.T) {
	instance := Instance{
		EncryptedSecureDetails: testArbitraryObjectJSON,
		SecureDetails:          SecureInstanceDetails{},
	}
	var err error
	instance, err = instance.decryptSecureDetails(noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testSecureInstanceDetails, instance.SecureDetails)
}
