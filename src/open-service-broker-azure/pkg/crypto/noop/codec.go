package noop

import "open-service-broker-azure/pkg/crypto"

type codec struct{}

// NewCodec returns a new no-op implementation of crypto.Codec
func NewCodec() crypto.Codec {
	return &codec{}
}

func (c *codec) Encrypt(plaintext []byte) ([]byte, error) {
	return plaintext, nil
}

func (c *codec) Decrypt(ciphertext []byte) ([]byte, error) {
	return ciphertext, nil
}
