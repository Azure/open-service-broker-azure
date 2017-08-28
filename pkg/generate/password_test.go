package generate

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPassword(t *testing.T) {
	overallRegex, err := regexp.Compile(`^[a-zA-Z\d]{16}$`)
	assert.Nil(t, err)
	lowerAlphaRegex, err := regexp.Compile(`[a-z]`)
	assert.Nil(t, err)
	upperAlphaRegex, err := regexp.Compile(`[A-Z]`)
	assert.Nil(t, err)
	numericRegex, err := regexp.Compile(`\d`)
	assert.Nil(t, err)
	for range [100]struct{}{} {
		password := NewPassword()
		assert.True(t, overallRegex.MatchString(password))
		assert.True(t, lowerAlphaRegex.MatchString(password))
		assert.True(t, upperAlphaRegex.MatchString(password))
		assert.True(t, numericRegex.MatchString(password))
	}
}
