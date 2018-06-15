package service

import (
	"encoding/json"

	"github.com/Azure/open-service-broker-azure/pkg/crypto"
)

// SecureString is a string that is seamlessly encrypted and decrypted when it
// is, respectively, marshaled or unmarshaled
type SecureString string

// MarshalJSON converts a SecureString to JSON, encrypting it in the process
func (s SecureString) MarshalJSON() ([]byte, error) {
	encryptedBytes, err := crypto.Encrypt([]byte(string(s)))
	if err != nil {
		return nil, err
	}
	return json.Marshal(encryptedBytes)
}

// UnmarshalJSON converts JSON to a SecureString, decrypting it in the process
func (s *SecureString) UnmarshalJSON(bytes []byte) error {
	var encryptedBytes []byte
	err := json.Unmarshal(bytes, &encryptedBytes)
	if err != nil {
		return err
	}
	decryptedBytes, err := crypto.Decrypt(encryptedBytes)
	if err != nil {
		return err
	}
	*s = SecureString(string(decryptedBytes))
	return nil
}
