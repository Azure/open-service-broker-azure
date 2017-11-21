package params

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddSingleValuedParameter(t *testing.T) {
	params := map[string]interface{}{}
	err := Add(params, "foo", "bar")
	assert.Nil(t, err)
	err = Add(params, "bat", "baz")
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"foo": "bar",
		"bat": "baz",
	}, params)

	// Test that it's a failure to now reuse a key in a mutli-valued context
	err = Add(params, "foo[0]", "bar")
	assert.NotNil(t, err)
}

func TestAddMultiValuedParameter(t *testing.T) {
	params := map[string]interface{}{}
	err := Add(params, "foo[0]", "bar")
	assert.Nil(t, err)
	err = Add(params, "foo[1]", "bat")
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"foo": []interface{}{"bar", "bat"},
	}, params)

	// Test that it's a failure to now reuse a key in a single-valued context
	err = Add(params, "foo", "bar")
	assert.NotNil(t, err)
}

func TestAddMultiValuedParametersOutOfOrder(t *testing.T) {
	params := map[string]interface{}{}
	err := Add(params, "foo[1]", "bar")
	assert.Nil(t, err)
	err = Add(params, "foo[0]", 5)
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"foo": []interface{}{5, "bar"},
	}, params)
}

func TestAddMultiValuedParametersWithGaps(t *testing.T) {
	params := map[string]interface{}{}
	err := Add(params, "foo[0]", "bar")
	assert.Nil(t, err)
	err = Add(params, "foo[2]", 5)
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"foo": []interface{}{"bar", nil, 5},
	}, params)
}
