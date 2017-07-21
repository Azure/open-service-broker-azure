package noop

import "github.com/Azure/azure-service-broker/pkg/crypto"

type codec struct{}

// NewCodec returns a new no-op implementation of crypto.Codec
func NewCodec() crypto.Codec {
	return &codec{}
}

func (c *codec) Encrypt(plaintext string) (string, error) {
	return plaintext, nil
}

func (c *codec) Decrypt(ciphertext string) (string, error) {
	return ciphertext, nil
}
