package fake

import "open-service-broker-azure/pkg/crypto"

// Codec is an implementation of crypto.Codec used to facilitate testing
type Codec struct {
	EncryptBehavior func(plaintext []byte) ([]byte, error)
	DecryptBehavior func(ciphertext []byte) ([]byte, error)
}

// NewCodec returns a new no-op implementation of crypto.Codec
func NewCodec() crypto.Codec {
	return &Codec{
		EncryptBehavior: defaultEncryptBehavior,
		DecryptBehavior: defaultDecryptBehavior,
	}
}

// Encrypt delegates encyption to an overridable function
func (c *Codec) Encrypt(plaintext []byte) ([]byte, error) {
	return c.EncryptBehavior(plaintext)
}

// Decrypt delegates decryption to an overridable function
func (c *Codec) Decrypt(ciphertext []byte) ([]byte, error) {
	return c.DecryptBehavior(ciphertext)
}

func defaultEncryptBehavior(plaintext []byte) ([]byte, error) {
	return plaintext, nil
}

func defaultDecryptBehavior(ciphertext []byte) ([]byte, error) {
	return ciphertext, nil
}
