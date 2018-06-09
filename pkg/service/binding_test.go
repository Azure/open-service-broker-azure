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
	testSecureBindingDetails = SecureBindingDetails{
		"foo": fooValue,
	}
	testBinding             Binding
	testBindingJSON         []byte
	bindingParametersSchema *InputParametersSchema
)

func init() {
	bindingID := "test-binding-id"
	instanceID := "test-instance-id"
	serviceID := "test-service-id"
	bindingParametersSchema = &InputParametersSchema{
		PropertySchemas: map[string]PropertySchema{
			"foo": &StringPropertySchema{},
		},
	}
	bindingParameters := &BindingParameters{
		Parameters: Parameters{
			Codec:  noopCodec,
			Schema: bindingParametersSchema,
			Data: map[string]interface{}{
				"foo": "bar",
			},
		},
	}
	bindingParametersJSONStr := []byte(`{"foo":"bar"}`)
	statusReason := "in-progress"
	bindingDetails := BindingDetails{
		"foo": "bat",
	}
	bindingDetailsJSONStr := `{"foo":"bat"}`
	secureBindingDetails := SecureBindingDetails{
		"foo": "bat",
	}
	encryptedSecureBindingDetails := []byte(`{"foo":"bat"}`)
	created, err := time.Parse(time.RFC3339, "2016-07-22T10:11:55-04:00")
	if err != nil {
		panic(err)
	}

	testBinding = Binding{
		BindingID:              bindingID,
		InstanceID:             instanceID,
		ServiceID:              serviceID,
		BindingParameters:      bindingParameters,
		Status:                 BindingStateBound,
		StatusReason:           statusReason,
		Details:                bindingDetails,
		EncryptedSecureDetails: encryptedSecureBindingDetails,
		SecureDetails:          secureBindingDetails,
		Created:                created,
	}

	b64EncryptedSecureBindingDetails := base64.StdEncoding.EncodeToString(
		encryptedSecureBindingDetails,
	)

	testBindingJSONStr := fmt.Sprintf(
		`{
			"bindingId":"%s",
			"instanceId":"%s",
			"serviceId":"%s",
			"bindingParameters":%s,
			"status":"%s",
			"statusReason":"%s",
			"details":%s,
			"secureDetails":"%s",
			"created":"%s"
		}`,
		bindingID,
		instanceID,
		serviceID,
		bindingParametersJSONStr,
		BindingStateBound,
		statusReason,
		bindingDetailsJSONStr,
		b64EncryptedSecureBindingDetails,
		created.Format(time.RFC3339),
	)
	testBindingJSONStr = strings.Replace(testBindingJSONStr, " ", "", -1)
	testBindingJSONStr = strings.Replace(testBindingJSONStr, "\n", "", -1)
	testBindingJSONStr = strings.Replace(testBindingJSONStr, "\t", "", -1)
	testBindingJSON = []byte(testBindingJSONStr)
}

func TestNewBindingFromJSON(t *testing.T) {
	binding, err := NewBindingFromJSON(
		testBindingJSON,
		noopCodec,
		bindingParametersSchema,
	)
	assert.Nil(t, err)
	assert.Equal(t, testBinding, binding)
}

func TestBindingToJSON(t *testing.T) {
	json, err := testBinding.ToJSON(noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testBindingJSON, json)
}

func TestEncryptSecureBindingDetails(t *testing.T) {
	binding := Binding{
		SecureDetails: testSecureBindingDetails,
	}
	var err error
	binding, err = binding.encryptSecureDetails(noopCodec)
	assert.Nil(t, err)
	assert.Equal(
		t,
		testArbitraryObjectJSON,
		binding.EncryptedSecureDetails,
	)
}

func TestDecryptSecureBindingDetails(t *testing.T) {
	binding := Binding{
		EncryptedSecureDetails: testArbitraryObjectJSON,
		SecureDetails:          SecureBindingDetails{},
	}
	var err error
	binding, err = binding.decryptSecureDetails(noopCodec)
	assert.Nil(t, err)
	assert.Equal(t, testSecureBindingDetails, binding.SecureDetails)
}
