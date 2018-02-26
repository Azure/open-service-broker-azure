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
			{
				Name:    "good rule",
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
	pp := &ServerProvisioningParameters{
		FirewallRules: []FirewallRule{
			{
				Name:    "good rule",
				StartIP: "192.168.86.1",
				EndIP:   "192.168.86.100",
			},
			{
				Name:    "good rule 2",
				StartIP: "192.168.86.101",
				EndIP:   "192.168.86.150",
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
			{
				StartIP: "192.168.86.1",
				EndIP:   "192.168.86.100",
			},
		},
	}
	error := sm.ValidateProvisioningParameters(pp, nil)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "name")
}

func TestValidateMissingEndFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := &ServerProvisioningParameters{
		FirewallRules: []FirewallRule{
			{
				Name:    "Test",
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
	pp := &ServerProvisioningParameters{
		FirewallRules: []FirewallRule{
			{
				Name:  "Test",
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
	pp := &ServerProvisioningParameters{
		FirewallRules: []FirewallRule{
			{
				Name:    "Test",
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
	pp := &ServerProvisioningParameters{
		FirewallRules: []FirewallRule{
			{
				Name:    "Test",
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
