package service

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalAndUnmarshalJSON(t *testing.T) {
	origStr := "foo"
	origSecureStr := SecureString(origStr)
	jsonBytes, err := json.Marshal(origSecureStr)
	assert.Nil(t, err)
	// Unmarshal into a plain string to assert encryption occurred during
	// marshaling
	var strFromJSON string
	err = json.Unmarshal(jsonBytes, &strFromJSON)
	assert.Nil(t, err)
	assert.NotEqual(t, origStr, strFromJSON)
	// Unmarshal into a secure string to assert decryption occurs during
	// unmarshalling
	var secureStr SecureString
	err = json.Unmarshal(jsonBytes, &secureStr)
	assert.Nil(t, err)
	assert.Equal(t, origSecureStr, secureStr)
}
