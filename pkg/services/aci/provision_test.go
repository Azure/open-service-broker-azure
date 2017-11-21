package aci

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateProvisioningParametersWithNoLocation(t *testing.T) {
	m, pp := getModuleAndValidTestProvisioningParameters(t)
	pp.Location = ""
	err := m.ValidateProvisioningParameters(pp)
	assert.NotNil(t, err)
}

func TestValidateProvisioningParametersWithInvalidLocation(t *testing.T) {
	m, pp := getModuleAndValidTestProvisioningParameters(t)
	pp.Location = "boguswest"
	err := m.ValidateProvisioningParameters(pp)
	assert.NotNil(t, err)
}

func TestValidateProvisioningParametersWithNoImageName(t *testing.T) {
	m, pp := getModuleAndValidTestProvisioningParameters(t)
	pp.ImageName = ""
	err := m.ValidateProvisioningParameters(pp)
	assert.NotNil(t, err)
}

func getModuleAndValidTestProvisioningParameters(
	t *testing.T,
) (*module, *ProvisioningParameters) {
	m, pp := &module{}, &ProvisioningParameters{
		Location:      "eastus",
		ResourceGroup: "test",
		ImageName:     "nginx:latest",
	}
	err := m.ValidateProvisioningParameters(pp)
	assert.Nil(t, err)
	return m, pp
}
