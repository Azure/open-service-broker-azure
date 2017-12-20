package sqldb

import (
	"testing"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/stretchr/testify/assert"
)

func TestValidateNoFirewallConfigAllInOne(t *testing.T) {

	sm := &allServiceManager{}

	pp := &ServerProvisioningParameters{}

	error := sm.ValidateProvisioningParameters(pp)
	assert.Nil(t, error)
}

func TestValidateNoFirewallConfigServerVMOnly(t *testing.T) {

	sm := &vmServiceManager{}

	pp := &ServerProvisioningParameters{}

	error := sm.ValidateProvisioningParameters(pp)
	assert.Nil(t, error)
}

func TestValidateGoodFirewallConfigAllInOne(t *testing.T) {

	sm := &allServiceManager{}

	pp := &ServerProvisioningParameters{
		FirewallIPStart: "192.168.86.1",
		FirewallIPEnd:   "192.168.86.100",
	}

	error := sm.ValidateProvisioningParameters(pp)
	assert.Nil(t, error)
}
func TestValidateGoodFirewallConfigServerVMOnly(t *testing.T) {

	sm := &vmServiceManager{}

	pp := &ServerProvisioningParameters{
		FirewallIPStart: "192.168.86.1",
		FirewallIPEnd:   "192.168.86.100",
	}

	error := sm.ValidateProvisioningParameters(pp)
	assert.Nil(t, error)
}

func TestValidateMissingEndFirewallConfigAllInOne(t *testing.T) {
	sm := &allServiceManager{}
	pp := &ServerProvisioningParameters{
		FirewallIPStart: "192.168.86.1",
	}
	error := sm.ValidateProvisioningParameters(pp)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallEndIPAddress")
}

func TestValidateMissingEndFirewallConfigServerVMOnly(t *testing.T) {
	sm := &vmServiceManager{}
	pp := &ServerProvisioningParameters{
		FirewallIPStart: "192.168.86.1",
	}
	error := sm.ValidateProvisioningParameters(pp)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallEndIPAddress")
}

func TestValidateMissingStartFirewallConfigAllInOne(t *testing.T) {
	sm := &allServiceManager{}
	pp := &ServerProvisioningParameters{
		FirewallIPEnd: "192.168.86.200",
	}
	error := sm.ValidateProvisioningParameters(pp)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallStartIPAddress")
}

func TestValidateMissingStartFirewallConfigServerVMOnly(t *testing.T) {
	sm := &vmServiceManager{}
	pp := &ServerProvisioningParameters{
		FirewallIPEnd: "192.168.86.200",
	}
	error := sm.ValidateProvisioningParameters(pp)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallStartIPAddress")
}
func TestValidateInvalidIPAllInOne(t *testing.T) {
	sm := &allServiceManager{}
	pp := &ServerProvisioningParameters{
		FirewallIPStart: "decafbad",
		FirewallIPEnd:   "192.168.86.200",
	}
	error := sm.ValidateProvisioningParameters(pp)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallStartIPAddress")
}
func TestValidateInvalidIPServerVMOnly(t *testing.T) {
	sm := &vmServiceManager{}
	pp := &ServerProvisioningParameters{
		FirewallIPStart: "decafbad",
		FirewallIPEnd:   "192.168.86.200",
	}
	error := sm.ValidateProvisioningParameters(pp)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallStartIPAddress")
}

func TestValidateIncompleteIPAllInOne(t *testing.T) {
	sm := &allServiceManager{}
	pp := &ServerProvisioningParameters{
		FirewallIPStart: "192.168.",
		FirewallIPEnd:   "192.168.86.200",
	}
	error := sm.ValidateProvisioningParameters(pp)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallStartIPAddress")
}
func TestValidateIncompleteIPServerVMOnly(t *testing.T) {
	sm := &vmServiceManager{}
	pp := &ServerProvisioningParameters{
		FirewallIPStart: "192.168.",
		FirewallIPEnd:   "192.168.86.200",
	}
	error := sm.ValidateProvisioningParameters(pp)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallStartIPAddress")
}
