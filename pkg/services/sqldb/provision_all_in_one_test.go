package sqldb

import (
	"testing"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/stretchr/testify/assert"
)

func TestValidateNoFirewallConfig(t *testing.T) {

	sm := &allInOneManager{}

	pp := &ServerProvisioningParams{}

	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.Nil(t, error)
}

func TestValidateGoodFirewallConfig(t *testing.T) {

	sm := &allInOneManager{}

	pp := &ServerProvisioningParams{
		FirewallRules: []FirewallRule{
			{
				Name:    "Goodrule",
				StartIP: "192.168.86.1",
				EndIP:   "192.168.86.100",
			},
		},
	}

	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.Nil(t, error)
}

func TestValidateMultipleGoodFirewallConfig(t *testing.T) {

	sm := &allInOneManager{}

	pp := &ServerProvisioningParams{
		FirewallRules: []FirewallRule{
			{
				Name:    "Goodrule",
				StartIP: "192.168.86.1",
				EndIP:   "192.168.86.100",
			},
			{
				Name:    "Goodrule2",
				StartIP: "192.168.86.101",
				EndIP:   "192.168.86.255",
			},
		},
	}

	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.Nil(t, error)
}

func TestValidateBadFirewallConfigMissingName(t *testing.T) {

	sm := &allInOneManager{}

	pp := &ServerProvisioningParams{
		FirewallRules: []FirewallRule{
			{
				StartIP: "192.168.86.1",
				EndIP:   "192.168.86.100",
			},
		},
	}

	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "name")
}

func TestValidateMissingEndFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := &ServerProvisioningParams{
		FirewallRules: []FirewallRule{
			{
				Name:    "BadRule",
				StartIP: "192.168.86.1",
			},
		},
	}
	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "endIPAddress")
}

func TestValidateMissingStartFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := &ServerProvisioningParams{
		FirewallRules: []FirewallRule{
			{
				Name:  "Badrule",
				EndIP: "192.168.86.200",
			},
		},
	}
	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "startIPAddress")
}

func TestValidateInvalidIP(t *testing.T) {
	sm := &allInOneManager{}
	pp := &ServerProvisioningParams{
		FirewallRules: []FirewallRule{
			{
				Name:    "BadRule",
				StartIP: "decafbad",
				EndIP:   "192.168.86.200",
			},
		},
	}
	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "startIPAddress")
}

func TestValidateIncompleteIP(t *testing.T) {
	sm := &allInOneManager{}
	pp := &ServerProvisioningParams{
		FirewallRules: []FirewallRule{
			{
				Name:    "Goodrule",
				StartIP: "192.168.",
				EndIP:   "192.168.86.200",
			},
		},
	}
	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "startIPAddress")
}
