package postgresqldb

import (
	"testing"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/stretchr/testify/assert"
)

func TestValidateNoFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := &AllInOneProvisioningParameters{}
	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.Nil(t, error)
}

func TestValidateGoodFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := &AllInOneProvisioningParameters{}
	pp.FirewallRules = []FirewallRule{
		{
			FirewallRuleName: "Good Rule",
			FirewallIPStart:  "192.168.86.1",
			FirewallIPEnd:    "192.168.86.100",
		},
	}

	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.Nil(t, error)
}

func TestValidateMissingFirewallRuleNameConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := &AllInOneProvisioningParameters{}
	pp.FirewallRules = []FirewallRule{
		{
			FirewallRuleName: "",
			FirewallIPStart:  "192.168.86.1",
			FirewallIPEnd:    "255.255.255.0",
		},
	}

	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallRuleName")
}
func TestValidateMissingEndFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := &AllInOneProvisioningParameters{}
	pp.FirewallRules = []FirewallRule{
		{
			FirewallRuleName: "Bad Rule",
			FirewallIPStart:  "192.168.86.1",
		},
	}

	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallEndIPAddress")
}

func TestValidateMissingStartFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := &AllInOneProvisioningParameters{}
	pp.FirewallRules = []FirewallRule{
		{
			FirewallRuleName: "BadRule",
			FirewallIPEnd:    "192.168.86.200",
		},
	}

	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallStartIPAddress")
}

func TestValidateInvalidIP(t *testing.T) {
	sm := &allInOneManager{}
	pp := &AllInOneProvisioningParameters{}
	pp.FirewallRules = []FirewallRule{
		{
			FirewallRuleName: "Bad Rule",
			FirewallIPStart:  "decafbad",
			FirewallIPEnd:    "192.168.86.200",
		},
	}
	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallStartIPAddress")
}

func TestValidateIncompleteIP(t *testing.T) {
	sm := &allInOneManager{}
	pp := &AllInOneProvisioningParameters{}
	pp.FirewallRules = []FirewallRule{
		{
			FirewallRuleName: "Bad Rule",
			FirewallIPStart:  "192.168.",
			FirewallIPEnd:    "192.168.86.200",
		},
	}

	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallStartIPAddress")
}
