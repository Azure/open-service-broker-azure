package service

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
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
			Schema: bindingParametersSchema,
			Data: map[string]interface{}{
				"foo": "bar",
			},
		},
	}
	bindingParametersJSONStr := []byte(`{"foo":"bar"}`)
	statusReason := "in-progress"
	bindingDetails := &arbitraryType{
		Foo: "bat",
	}
	bindingDetailsJSONStr := `{"foo":"bat"}`
	created, err := time.Parse(time.RFC3339, "2016-07-22T10:11:55-04:00")
	if err != nil {
		panic(err)
	}

	testBinding = Binding{
		BindingID:         bindingID,
		InstanceID:        instanceID,
		ServiceID:         serviceID,
		BindingParameters: bindingParameters,
		Status:            BindingStateBound,
		StatusReason:      statusReason,
		Details:           bindingDetails,
		Created:           created,
	}

	testBindingJSONStr := fmt.Sprintf(
		`{
			"bindingId":"%s",
			"instanceId":"%s",
			"serviceId":"%s",
			"bindingParameters":%s,
			"status":"%s",
			"statusReason":"%s",
			"details":%s,
			"created":"%s"
		}`,
		bindingID,
		instanceID,
		serviceID,
		bindingParametersJSONStr,
		BindingStateBound,
		statusReason,
		bindingDetailsJSONStr,
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
		&arbitraryType{},
		bindingParametersSchema,
	)
	assert.Nil(t, err)
	assert.Equal(t, testBinding, binding)
}

func TestBindingToJSON(t *testing.T) {
	json, err := testBinding.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, testBindingJSON, json)
}
