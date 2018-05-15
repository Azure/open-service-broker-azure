package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptySlice(t *testing.T) {
	var slice []string
	assert.True(t, IsEmpty(slice))
}

func TestNonEmptySlice(t *testing.T) {
	slice := []string{"hello", "blah"}
	assert.False(t, IsEmpty(slice))
}

func TestEmptyString(t *testing.T) {
	value := ""
	assert.True(t, IsEmpty(value))
}

func TestNonEmptyString(t *testing.T) {
	value := "notempty"
	assert.False(t, IsEmpty(value))
}

func TestEmptyIntPointer(t *testing.T) {
	var intPtr *int
	assert.True(t, IsEmpty(intPtr))
}

func TestNonEmptyIntPointer(t *testing.T) {
	var intPtr *int
	val := 5
	intPtr = &val
	assert.False(t, IsEmpty(intPtr))
}

func TestEmptyMap(t *testing.T) {
	val := make(map[string]interface{})
	assert.True(t, IsEmpty(val))
}

func TestNonEmptyMap(t *testing.T) {
	val := make(map[string]interface{})
	val["key"] = "value"
	assert.False(t, IsEmpty(val))
}
