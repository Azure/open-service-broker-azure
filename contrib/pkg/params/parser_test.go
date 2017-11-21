package params

import "testing"
import "github.com/stretchr/testify/assert"

func TestParseParam(t *testing.T) {
	// A param is invalid if it contains no "="
	_, _, err := Parse("foo")
	assert.NotNil(t, err)

	// Test a normal scenario
	key, val, err := Parse("foo=bar")
	assert.Nil(t, err)
	assert.Equal(t, "foo", key)
	assert.Equal(t, "bar", val)

	// Test that complex (indexed) keys are ok
	key, val, err = Parse("foo[0]=bar")
	assert.Nil(t, err)
	assert.Equal(t, "foo[0]", key)
	assert.Equal(t, "bar", val)

	// Test that multiple "=" are not a problem-- i.e. values can contain "="
	// ("=" is a base64 character and therefore somewhat common!)
	key, val, err = Parse("foo=bar==")
	assert.Nil(t, err)
	assert.Equal(t, "foo", key)
	assert.Equal(t, "bar==", val)
}
