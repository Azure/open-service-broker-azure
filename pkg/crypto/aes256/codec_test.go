package aes256

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodecEncryptAndDecrypt(t *testing.T) {
	c, err := NewCodec([]byte("AES256Key-32Characters1234567890"))
	assert.Nil(t, err)
	initialPlaintext := []byte("foo")
	ciphertext, err := c.Encrypt(initialPlaintext)
	assert.Nil(t, err)
	assert.NotEqual(t, initialPlaintext, ciphertext)
	plaintext, err := c.Decrypt(ciphertext)
	assert.Nil(t, err)
	assert.Equal(t, initialPlaintext, plaintext)
}
