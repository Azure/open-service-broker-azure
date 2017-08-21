package generate

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPassword(t *testing.T) {
	regex, err := regexp.Compile(`^[a-zA-Z\d]{16}$`)
	assert.Nil(t, err)
	for range [100]struct{}{} {
		password := NewPassword()
		assert.True(t, regex.MatchString(password))
	}
}
