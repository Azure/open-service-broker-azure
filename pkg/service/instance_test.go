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
	updatingParameters := ProvisioningParameters{
		"foo": "bat",
	}
	updatingParametersJSONStr := []byte(`{"foo":"bat"}`)
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
		InstanceID:             instanceID,
		Alias:                  alias,
		ServiceID:              serviceID,
		PlanID:                 planID,
		ProvisioningParameters: provisioningParameters,
		UpdatingParameters:     updatingParameters,
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
			"updatingParameters":%s,
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
		updatingParametersJSONStr,
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
