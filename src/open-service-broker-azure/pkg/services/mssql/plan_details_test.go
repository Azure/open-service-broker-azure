package mssql

import (
	"testing"

	"open-service-broker-azure/pkg/service"
	"github.com/stretchr/testify/assert"
)

func TestValidateInvalidIP(t *testing.T) {
	err := ipValidator("", "decafbad")
	assert.NotNil(t, err)
	err = ipValidator("", "192.168.")
	assert.NotNil(t, err)
	_, ok := err.(*service.ValidationError)
	assert.True(t, ok)
}

func TestValidateValidIP(t *testing.T) {
	err := ipValidator("", "192.168.1.100")
	assert.Nil(t, err)
	err = ipValidator("", "10.0.1.101")
	assert.Nil(t, err)
}

func TestValidateFirewallRuleWithEndIPBeforeStartIP(t *testing.T) {
	fr := map[string]interface{}{
		"name":           "Good Rule",
		"startIPAddress": "192.168.86.100",
		"endIPAddress":   "192.168.86.1",
	}
	err := firewallRuleValidator("", fr)
	assert.NotNil(t, err)
	_, ok := err.(*service.ValidationError)
	assert.True(t, ok)
}

func TestValidateValidFirewallRule(t *testing.T) {
	fr := map[string]interface{}{
		"name":           "Good Rule",
		"startIPAddress": "192.168.86.1",
		"endIPAddress":   "192.168.86.100",
	}
	err := firewallRuleValidator("", fr)
	assert.Nil(t, err)
}
