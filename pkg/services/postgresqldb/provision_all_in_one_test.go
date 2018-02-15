package postgresqldb

import (
	"testing"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/stretchr/testify/assert"
)

func TestValidateNoFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := &AllInOneProvisioningParameters{}
	error := sm.ValidateProvisioningParameters(pp)
	assert.Nil(t, error)
}

func TestValidateGoodFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := &AllInOneProvisioningParameters{}
	pp.FirewallIPStart = "192.168.86.1"
	pp.FirewallIPEnd = "192.168.86.100"

	error := sm.ValidateProvisioningParameters(pp)
	assert.Nil(t, error)
}

func TestValidateMissingEndFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := &AllInOneProvisioningParameters{}
	pp.FirewallIPStart = "192.168.86.1"
	error := sm.ValidateProvisioningParameters(pp)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallEndIPAddress")
}

func TestValidateMissingStartFirewallConfig(t *testing.T) {
	sm := &allInOneManager{}
	pp := &AllInOneProvisioningParameters{}
	pp.FirewallIPEnd = "192.168.86.200"

	error := sm.ValidateProvisioningParameters(pp)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallStartIPAddress")
}

func TestValidateInvalidIP(t *testing.T) {
	sm := &allInOneManager{}
	pp := &AllInOneProvisioningParameters{}
	pp.FirewallIPStart = "decafbad"
	pp.FirewallIPEnd = "192.168.86.200"
	error := sm.ValidateProvisioningParameters(pp)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallStartIPAddress")
}

func TestValidateIncompleteIP(t *testing.T) {
	sm := &allInOneManager{}
	pp := &AllInOneProvisioningParameters{}
	pp.FirewallIPStart = "192.168."
	pp.FirewallIPEnd = "192.168.86.200"
	error := sm.ValidateProvisioningParameters(pp)
	assert.NotNil(t, error)
	v, ok := error.(*service.ValidationError)
	assert.True(t, ok)
	assert.Equal(t, v.Field, "firewallStartIPAddress")
}
