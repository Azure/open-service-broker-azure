package aes256

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/Azure/azure-service-broker/pkg/crypto"
)

const nonceLength = 12

type codec struct {
	aesgcm cipher.AEAD
}

// NewCodec returns a new aes256-based implementation of crypto.Codec
func NewCodec(key []byte) (crypto.Codec, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating cipher: %s", err)
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error creating cipher: %s", err)
	}
	return &codec{
		aesgcm: aesgcm,
	}, nil
}

func (c *codec) Encrypt(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, nonceLength)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("error generating nonce: %s", err)
	}
	ciphertext := c.aesgcm.Seal(nil, nonce, plaintext, nil)
	// Return the ciphertext prefixed with the nonce-- this consolidates both
	// into a single value so that anyone who has encrypted using this scheme
	// isn't burdened with schlepping / storing the nonce in addition to the
	// ciphertext. The Decrypt() function simply possesses the intelligence to
	// split the nonce from the rest of the ciphertext before proceeding with
	// decryption.
	return append(nonce, ciphertext...), nil
}

func (c *codec) Decrypt(ciphertext []byte) ([]byte, error) {
	nonce := ciphertext[:nonceLength]
	ciphertext = ciphertext[nonceLength:]
	plaintext, err := c.aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("error decrypting ciphertext: %s", err)
	}
	return plaintext, nil
}
