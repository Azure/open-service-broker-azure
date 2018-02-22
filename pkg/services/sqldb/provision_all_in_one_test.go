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
		FirewallIPStart: "192.168.86.1",
		FirewallIPEnd:   "192.168.86.100",
	}

	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.Nil(t, error)
}

func TestValidateMissingEndFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := &ServerProvisioningParams{
		FirewallIPStart: "192.168.86.1",
	}
	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallEndIPAddress")
}

func TestValidateMissingStartFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := &ServerProvisioningParams{
		FirewallIPEnd: "192.168.86.200",
	}
	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallStartIPAddress")
}

func TestValidateInvalidIP(t *testing.T) {
	sm := &allInOneManager{}
	pp := &ServerProvisioningParams{
		FirewallIPStart: "decafbad",
		FirewallIPEnd:   "192.168.86.200",
	}
	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallStartIPAddress")
}

func TestValidateIncompleteIP(t *testing.T) {
	sm := &allInOneManager{}
	pp := &ServerProvisioningParams{
		FirewallIPStart: "192.168.",
		FirewallIPEnd:   "192.168.86.200",
	}
	error := sm.ValidateProvisioningParameters(pp, nil)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallStartIPAddress")
}
