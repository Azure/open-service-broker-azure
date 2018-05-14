package aci

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func memoryValidatorTest(t *testing.T) {
	// Some invalid cases
	err := memoryValidator("", 0.05)
	assert.NotNil(t, err)
	err = memoryValidator("", 1.25)
	assert.NotNil(t, err)
	// Some valid cases
	err = memoryValidator("", 0.1)
	assert.Nil(t, err)
	err = memoryValidator("", 0.5)
	assert.Nil(t, err)
	err = memoryValidator("", 2.0)
	assert.Nil(t, err)
}
