package postgresql

import (
	"testing"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/stretchr/testify/assert"
)

func TestValidateStorageIncreases(t *testing.T) {
	sm := &dbmsManager{}
	pp := service.Parameters{
		"storage": 10,
	}

	up := service.Parameters{
		"storage": 20,
	}

	plan := service.NewPlan(
		createBasicPlan("73191861-04b3-4d0b-a29b-429eb15a83d4", false),
	)

	instance := service.Instance{
		Plan: plan,
		ProvisioningParameters: pp,
		UpdatingParameters:     up,
	}

	err := sm.ValidateUpdatingParameters(instance)
	assert.Nil(t, err)
}

func TestValidateStorageDecreaseFails(t *testing.T) {
	sm := &dbmsManager{}
	pp := service.Parameters{
		"storage": 20,
	}

	up := service.Parameters{
		"storage": 10,
	}

	plan := service.NewPlan(
		createBasicPlan("73191861-04b3-4d0b-a29b-429eb15a83d4", false),
	)

	instance := service.Instance{
		Plan: plan,
		ProvisioningParameters: pp,
		UpdatingParameters:     up,
	}

	err := sm.ValidateUpdatingParameters(instance)
	assert.NotNil(t, err)
	v, ok := err.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "storage", v.Field)
}
