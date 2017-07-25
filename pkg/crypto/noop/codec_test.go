package noop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCodec = NewCodec()

func TestCodecEncrypt(t *testing.T) {
	plaintext := "foo"
	ciphertext, err := testCodec.Encrypt(plaintext)
	assert.Nil(t, err)
	assert.Equal(t, plaintext, ciphertext)
}

func TestCodecDecrypt(t *testing.T) {
	ciphertext := "foo"
	plaintext, err := testCodec.Decrypt(ciphertext)
	assert.Nil(t, err)
	assert.Equal(t, ciphertext, plaintext)
}
