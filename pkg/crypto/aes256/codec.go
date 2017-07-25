package aes256

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/Azure/azure-service-broker/pkg/crypto"
)

type codec struct {
	aesgcm cipher.AEAD
}

// NewCodec returns a new aes256-based implementation of crypto.Codec
func NewCodec(key string) (crypto.Codec, error) {
	block, err := aes.NewCipher([]byte(key))
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

func (c *codec) Encrypt(plaintext string) (string, error) {
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("error generating nonce: %s", err)
	}
	b64NonceStr := base64.StdEncoding.EncodeToString(nonce)
	ciphertextBytes := c.aesgcm.Seal(nil, nonce, []byte(plaintext), nil)
	b64Ciphertext := base64.StdEncoding.EncodeToString(ciphertextBytes)
	// Return the ciphertext prefixed with the nonce-- this consolidates both
	// into a single value so that anyone who has encrypted using this scheme
	// isn't burdened with schlepping / storing the nonce in addition to the
	// ciphertext. The Decrypt() function simply posseses the intelligence to
	// split the nonce from the rest of the ciphertext before procedding with
	// decryption.
	ret := fmt.Sprintf("%s:%s", b64NonceStr, b64Ciphertext)
	return ret, nil
}

func (c *codec) Decrypt(ciphertext string) (string, error) {
	// We prefix ciphertext with the nonce-- split these
	tokens := strings.SplitN(ciphertext, ":", 2)
	if len(tokens) != 2 {
		return "", fmt.Errorf("invalid ciphertext: %s", ciphertext)
	}
	nonce, err := base64.StdEncoding.DecodeString(tokens[0])
	if err != nil {
		return "", fmt.Errorf("error decoding nonce: %s", err)
	}
	ciphertextBytes, err := base64.StdEncoding.DecodeString(tokens[1])
	if err != nil {
		return "", fmt.Errorf("error decoding ciphertext: %s", err)
	}
	plaintext, err := c.aesgcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", fmt.Errorf("error decrypting ciphertext: %s", err)
	}
	return string(plaintext), nil
}
