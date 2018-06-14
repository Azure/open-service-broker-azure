package service

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testInstance                 Instance
	testInstanceJSON             []byte
	provisioningParametersSchema *InputParametersSchema
)

func init() {
	instanceID := "test-instance-id"
	alias := "test-alias"
	serviceID := "test-service-id"
	planID := "test-plan-id"
	parentAlias := "test-parent-alias"
	provisioningParametersSchema = &InputParametersSchema{
		PropertySchemas: map[string]PropertySchema{
			"foo": &StringPropertySchema{},
		},
	}
	provisioningParameters := &ProvisioningParameters{
		Parameters: Parameters{
			Schema: provisioningParametersSchema,
			Data: map[string]interface{}{
				"foo": "bar",
			},
		},
	}
	provisioningParametersJSONStr := []byte(`{"foo":"bar"}`)
	updatingParameters := &ProvisioningParameters{
		Parameters: Parameters{
			Schema: provisioningParametersSchema,
			Data: map[string]interface{}{
				"foo": "bat",
			},
		},
	}
	updatingParametersJSONStr := []byte(`{"foo":"bat"}`)
	statusReason := "in-progress"
	details := &arbitraryType{
		Foo: "baz",
	}
	detailsJSONStr := `{"foo":"baz"}`
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
		ParentAlias:            parentAlias,
		Details:                details,
		Created:                created,
	}

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
			"parentAlias":"%s",
			"details":%s,
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
		parentAlias,
		detailsJSONStr,
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
		&arbitraryType{},
		provisioningParametersSchema,
	)
	assert.Nil(t, err)
	assert.Equal(t, testInstance, instance)
}

func TestInstanceToJSON(t *testing.T) {
	json, err := testInstance.ToJSON()
	assert.Nil(t, err)
	assert.Equal(t, string(testInstanceJSON), string(json))
}
