package postgresql

import (
	"testing"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/stretchr/testify/assert"
)

func TestValidateStorageIncreases(t *testing.T) {
	sm := &dbmsManager{}

	schema := &service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{
			"storage": &service.IntPropertySchema{},
		},
	}
	instance := service.Instance{
		ProvisioningParameters: service.Parameters{
			Schema: schema,
			Data: map[string]interface{}{
				"storage": 10,
			},
		},
		UpdatingParameters: service.Parameters{
			Schema: schema,
			Data: map[string]interface{}{
				"storage": 20,
			},
		},
	}

	err := sm.ValidateUpdatingParameters(instance)
	assert.Nil(t, err)
}

func TestValidateStorageDecreaseFails(t *testing.T) {
	sm := &dbmsManager{}

	schema := &service.InputParametersSchema{
		PropertySchemas: map[string]service.PropertySchema{
			"storage": &service.IntPropertySchema{},
		},
	}
	instance := service.Instance{
		ProvisioningParameters: service.Parameters{
			Schema: schema,
			Data: map[string]interface{}{
				"storage": 20,
			},
		},
		UpdatingParameters: service.Parameters{
			Schema: schema,
			Data: map[string]interface{}{
				"storage": 10,
			},
		},
	}

	err := sm.ValidateUpdatingParameters(instance)
	assert.NotNil(t, err)
	v, ok := err.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "storage", v.Field)
}
