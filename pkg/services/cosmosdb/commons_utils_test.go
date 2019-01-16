package cosmosdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateReadLocations(t *testing.T) {
	assert.NotNil(t, validateReadLocations("", []string{"ukwest", "no-existing"}))
	assert.NotNil(t, validateReadLocations("", []string{"ukwest", "ukwest"}))

	assert.Nil(t, validateReadLocations("", []string{"ukwest", "eastasia"}))
	assert.Nil(t, validateReadLocations("", []string{}))
}
