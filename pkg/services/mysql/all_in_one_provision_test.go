package mysql

import (
	"testing"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/stretchr/testify/assert"
)

func TestValidateNoFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := service.ProvisioningParameters{}
	plan := service.NewPlan(
		createBasicPlan("f7a86f81-0384-4999-a404-537a564abb62"),
	)
	err := sm.ValidateProvisioningParameters(plan, pp, nil)
	assert.Nil(t, err)
}

func TestValidateGoodFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := service.ProvisioningParameters{
		"firewallRules": []map[string]string{
			{
				"name":           "Good Rule",
				"startIPAddress": "192.168.86.1",
				"endIPAddress":   "192.168.86.100",
			},
		},
	}
	plan := service.NewPlan(
		createBasicPlan("f7a86f81-0384-4999-a404-537a564abb62"),
	)
	err := sm.ValidateProvisioningParameters(plan, pp, nil)
	assert.Nil(t, err)
}

func TestValidateMultipleGoodFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := service.ProvisioningParameters{
		"firewallRules": []map[string]string{
			{
				"name":           "Good Rule",
				"startIPAddress": "192.168.86.1",
				"endIPAddress":   "192.168.86.100",
			},
			{
				"name":           "Good Rule 2",
				"startIPAddress": "192.168.86.101",
				"endIPAddress":   "192.168.86.150",
			},
		},
	}
	plan := service.NewPlan(
		createBasicPlan("f7a86f81-0384-4999-a404-537a564abb62"),
	)
	err := sm.ValidateProvisioningParameters(plan, pp, nil)
	assert.Nil(t, err)
}

func TestValidateMissingFirewallRuleName(t *testing.T) {
	sm := &allInOneManager{}
	pp := service.ProvisioningParameters{
		"firewallRules": []map[string]string{
			{
				"startIPAddress": "192.168.86.1",
				"endIPAddress":   "192.168.86.100",
			},
		},
	}
	plan := service.NewPlan(
		createBasicPlan("f7a86f81-0384-4999-a404-537a564abb62"),
	)
	err := sm.ValidateProvisioningParameters(plan, pp, nil)
	v, ok := err.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "name", v.Field)
}

func TestValidateMissingEndFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := service.ProvisioningParameters{
		"firewallRules": []map[string]string{
			{
				"name":           "Bad Rule",
				"startIPAddress": "192.168.86.1",
			},
		},
	}
	plan := service.NewPlan(
		createBasicPlan("f7a86f81-0384-4999-a404-537a564abb62"),
	)
	err := sm.ValidateProvisioningParameters(plan, pp, nil)
	assert.NotNil(t, err)
	v, ok := err.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "endIPAddress", v.Field)
}

func TestValidateMissingStartFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := service.ProvisioningParameters{
		"firewallRules": []map[string]string{
			{
				"name":         "Bad Rule",
				"endIPAddress": "192.168.86.100",
			},
		},
	}
	plan := service.NewPlan(
		createBasicPlan("f7a86f81-0384-4999-a404-537a564abb62"),
	)
	err := sm.ValidateProvisioningParameters(plan, pp, nil)
	assert.NotNil(t, err)
	v, ok := err.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "startIPAddress", v.Field)
}

func TestValidateInvalidIP(t *testing.T) {
	sm := &allInOneManager{}
	pp := service.ProvisioningParameters{
		"firewallRules": []map[string]string{
			{
				"name":           "Bad Rule",
				"startIPAddress": "decafbad",
				"endIPAddress":   "192.168.86.100",
			},
		},
	}
	plan := service.NewPlan(
		createBasicPlan("f7a86f81-0384-4999-a404-537a564abb62"),
	)
	err := sm.ValidateProvisioningParameters(plan, pp, nil)
	assert.NotNil(t, err)
	v, ok := err.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "startIPAddress", v.Field)
}

func TestValidateIncompleteIP(t *testing.T) {
	sm := &allInOneManager{}
	pp := service.ProvisioningParameters{
		"firewallRules": []map[string]string{
			{
				"name":           "Bad Rule",
				"startIPAddress": "192.168",
				"endIPAddress":   "192.168.86.100",
			},
		},
	}
	plan := service.NewPlan(
		createBasicPlan("f7a86f81-0384-4999-a404-537a564abb62"),
	)
	err := sm.ValidateProvisioningParameters(plan, pp, nil)
	assert.NotNil(t, err)
	v, ok := err.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "startIPAddress", v.Field)
}

func TestValidateHardwareVersionIncompatible(t *testing.T) {
	provisionSchema := planSchema{
		allowedHardware:         []string{"gen5"},
		defaultHardware:         "gen5",
		allowedCores:            []int64{2, 4, 8, 16},
		defaultCores:            2,
		maxStorage:              2048,
		minStorage:              5,
		defaultStorage:          10,
		allowedBackupRedundancy: []string{"local", "geo"},
		minBackupRetention:      7,
		maxBackupRetention:      35,
		defaultBackupRetention:  7,
		tier: "MO",
	}
	extendedPlanData := map[string]interface{}{
		"provisionSchema": provisionSchema,
		"tier":            "MemoryOptimized",
	}

	plan := service.NewPlan(&service.PlanProperties{
		ID:          "f7a86f81-0384-4999-a404-537a564abb62",
		Name:        "somePlan",
		Description: "somePlan",
		Free:        false,
		Extended:    extendedPlanData,
		Metadata: &service.ServicePlanMetadata{
			DisplayName: "somePlan",
			Bullets:     []string{"Testable"},
		},
		Schemas: service.PlanSchemas{
			ServiceInstances: service.InstanceSchemas{
				ProvisioningParametersSchema: generateDBMSPlanSchema(provisionSchema),
			},
		},
	})

	sm := &allInOneManager{}
	pp := service.ProvisioningParameters{
		"hardwareFamily": "gen4",
		"firewallRules": []map[string]string{
			{
				"name":           "Good Rule",
				"startIPAddress": "192.168.86.1",
				"endIPAddress":   "192.168.86.100",
			},
		},
	}
	err := sm.ValidateProvisioningParameters(plan, pp, nil)
	assert.NotNil(t, err)
	v, ok := err.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, "hardwareFamily", v.Field)

}
