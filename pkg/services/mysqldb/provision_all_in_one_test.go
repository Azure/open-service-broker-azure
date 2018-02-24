package mysqldb

import (
	"testing"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/stretchr/testify/assert"
)

func TestValidateNoFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := &ServerProvisioningParameters{}
	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.Nil(t, error)
}

func TestValidateGoodFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := &ServerProvisioningParameters{
		FirewallRules: []FirewallRule{
			FirewallRule{
				FirewallRuleName: "good rule",
				FirewallIPStart:  "192.168.86.1",
				FirewallIPEnd:    "192.168.86.100",
			},
		},
	}
	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.Nil(t, error)
}

func TestValidateMultipleGoodFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := &ServerProvisioningParameters{
		FirewallRules: []FirewallRule{
			FirewallRule{
				FirewallRuleName: "good rule",
				FirewallIPStart:  "192.168.86.1",
				FirewallIPEnd:    "192.168.86.100",
			},
			FirewallRule{
				FirewallRuleName: "good rule 2",
				FirewallIPStart:  "192.168.86.101",
				FirewallIPEnd:    "192.168.86.150",
			},
		},
	}
	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.Nil(t, error)
}

func TestValidateMissingFirewallRuleName(t *testing.T) {
	sm := &allInOneManager{}
	pp := &ServerProvisioningParameters{
		FirewallRules: []FirewallRule{
			FirewallRule{
				FirewallIPStart: "192.168.86.1",
				FirewallIPEnd:   "192.168.86.100",
			},
		},
	}
	error := sm.ValidateProvisioningParameters(pp, nil)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallRuleName")
}

func TestValidateMissingEndFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := &ServerProvisioningParameters{
		FirewallRules: []FirewallRule{
			FirewallRule{
				FirewallRuleName: "Test",
				FirewallIPStart:  "192.168.86.1",
			},
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
	pp := &ServerProvisioningParameters{
		FirewallRules: []FirewallRule{
			FirewallRule{
				FirewallRuleName: "Test",
				FirewallIPEnd:    "192.168.86.200",
			},
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
	pp := &ServerProvisioningParameters{
		FirewallRules: []FirewallRule{
			FirewallRule{
				FirewallRuleName: "Test",
				FirewallIPStart:  "decafbad",
				FirewallIPEnd:    "192.168.86.200",
			},
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
	pp := &ServerProvisioningParameters{
		FirewallRules: []FirewallRule{
			FirewallRule{
				FirewallRuleName: "Test",
				FirewallIPStart:  "192.168.",
				FirewallIPEnd:    "192.168.86.200",
			},
		},
	}
	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallStartIPAddress")
}
