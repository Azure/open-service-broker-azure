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
			Name:    "Good Rule",
			StartIP: "192.168.86.1",
			EndIP:   "192.168.86.100",
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
			Name:    "",
			StartIP: "192.168.86.1",
			EndIP:   "255.255.255.0",
		},
	}

	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "ruleName")
}
func TestValidateMissingEndFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := &AllInOneProvisioningParameters{}
	pp.FirewallRules = []FirewallRule{
		{
			Name:    "Bad Rule",
			StartIP: "192.168.86.1",
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
	pp := &AllInOneProvisioningParameters{}
	pp.FirewallRules = []FirewallRule{
		{
			Name:  "BadRule",
			EndIP: "192.168.86.200",
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
	pp := &AllInOneProvisioningParameters{}
	pp.FirewallRules = []FirewallRule{
		{
			Name:    "Bad Rule",
			StartIP: "decafbad",
			EndIP:   "192.168.86.200",
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
	pp := &AllInOneProvisioningParameters{}
	pp.FirewallRules = []FirewallRule{
		{
			Name:    "Bad Rule",
			StartIP: "192.168.",
			EndIP:   "192.168.86.200",
		},
	}

	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "startIPAddress")
}
