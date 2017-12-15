package mysqldb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateGoodFirewallConfig(t *testing.T) {

	sm := &serviceManager{}

	pp := &ProvisioningParameters{
		FirewallIPStart: "192.168.86.1",
		FirewallIPEnd:   "192.168.86.100",
	}

	error := sm.ValidateProvisioningParameters(pp)
	assert.Nil(t, error)
}

func TestValidateMissingEndFirewallConfig(t *testing.T) {

	sm := &serviceManager{}

	pp := &ProvisioningParameters{
		FirewallIPStart: "192.168.86.1",
	}

	error := sm.ValidateProvisioningParameters(pp)
	assert.NotNil(t, error)
}

func TestValidateInvalidIP(t *testing.T) {

	sm := &serviceManager{}

	pp := &ProvisioningParameters{
		FirewallIPStart: "decafbad",
	}

	error := sm.ValidateProvisioningParameters(pp)
	assert.NotNil(t, error)
}
