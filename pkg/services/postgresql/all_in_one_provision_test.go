package postgresql

import (
	"testing"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/stretchr/testify/assert"
)

func TestValidateNoFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := service.ProvisioningParameters{}
	err := sm.ValidateProvisioningParameters(pp, nil)
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
	err := sm.ValidateProvisioningParameters(pp, nil)
	assert.Nil(t, err)
}

func TestValidateMissingFirewallRuleNameConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := service.ProvisioningParameters{
		"firewallRules": []map[string]string{
			{
				"name":           "",
				"startIPAddress": "192.168.86.1",
				"endIPAddress":   "255.255.255.0",
			},
		},
	}
	err := sm.ValidateProvisioningParameters(pp, nil)
	assert.NotNil(t, err)
	v, ok := err.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "ruleName")
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
	err := sm.ValidateProvisioningParameters(pp, nil)
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
				"endIPAddress": "255.255.255.0",
			},
		},
	}
	err := sm.ValidateProvisioningParameters(pp, nil)
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
				"endIPAddress":   "255.255.255.0",
			},
		},
	}
	err := sm.ValidateProvisioningParameters(pp, nil)
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
				"startIPAddress": "192.168.",
				"endIPAddress":   "255.255.255.0",
			},
		},
	}
	err := sm.ValidateProvisioningParameters(pp, nil)
	assert.NotNil(t, err)
	v, ok := err.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "startIPAddress")
}
