package cosmosdb

import (
	"testing"

	"github.com/Azure/open-service-broker-azure/pkg/service"
	"github.com/stretchr/testify/assert"
)

func TestValidateInvalidIPRange(t *testing.T) {
	err := ipRangeValidator("", "decafbad")
	assert.NotNil(t, err)
	_, ok := err.(*service.ValidationError)
	assert.True(t, ok)
	err = ipRangeValidator("", "192.168.")
	assert.NotNil(t, err)
	_, ok = err.(*service.ValidationError)
	assert.True(t, ok)
}

func TestValidateValidIPRanges(t *testing.T) {
	err := ipRangeValidator("", "192.168.1.100")
	assert.Nil(t, err)
	err = ipRangeValidator("", "10.0.0.0/16")
	assert.Nil(t, err)
}
