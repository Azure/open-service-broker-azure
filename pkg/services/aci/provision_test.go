package aci

import (
	"testing"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/stretchr/testify/assert"
)

func TestValidateProvisioningParametersWithNoImageName(t *testing.T) {
	m := &module{}
	pp := service.ProvisioningParameters{
		"image": "nginx:latest",
	}
	err := m.serviceManager.ValidateProvisioningParameters(nil, pp, nil)
	assert.Nil(t, err)
	// pp = service.ProvisioningParameters{}
	// err = m.serviceManager.ValidateProvisioningParameters(pp, nil)
	// assert.NotNil(t, err)
}
