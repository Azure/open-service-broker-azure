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
	assert.Equal(t, v.Field, "name")
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
	assert.Equal(t, v.Field, "endIPAddress")
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
	assert.Equal(t, v.Field, "startIPAddress")
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
	assert.Equal(t, v.Field, "startIPAddress")
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
	assert.Equal(t, v.Field, "startIPAddress")
}
