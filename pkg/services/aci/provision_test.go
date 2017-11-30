package aci

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateProvisioningParametersWithNoImageName(t *testing.T) {
	m := &module{}
	pp := &ProvisioningParameters{
		ImageName: "nginx:latest",
	}
	err := m.serviceManager.ValidateProvisioningParameters(pp)
	assert.Nil(t, err)
	pp.ImageName = ""
	err = m.serviceManager.ValidateProvisioningParameters(pp)
	assert.NotNil(t, err)
}
