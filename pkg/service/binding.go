package service

import (
	"encoding/json"

	"github.com/Azure/azure-service-broker/pkg/crypto"
)

// Binding represents a biding to a service
type Binding struct {
	BindingID                  string `json:"bindingId"`
	InstanceID                 string `json:"instanceId"`
	EncryptedBindingParameters string `json:"bindingParameters"`
	Status                     string `json:"status"`
	StatusReason               string `json:"statusReason"`
	EncryptedBindingContext    string `json:"bindingContext"`
	EncryptedCredentials       string `json:"credentials"`
}

// NewBindingFromJSONString returns a new Binding unmarshalled from the
// provided JSON string
func NewBindingFromJSONString(jsonStr string) (*Binding, error) {
	binding := &Binding{}
	err := json.Unmarshal([]byte(jsonStr), binding)
	if err != nil {
		return nil, err
	}
	return binding, nil
}

// ToJSONString returns a string containing a JSON representation of the
// instance
func (b *Binding) ToJSONString() (string, error) {
	bytes, err := json.Marshal(b)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// SetBindingParameters marshals the provided bindingParameters object, encrypts
// the result, and stores it in the EncryptedBindingParameters field
func (b *Binding) SetBindingParameters(
	params interface{},
	codec crypto.Codec,
) error {
	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return err
	}
	ciphertext, err := codec.Encrypt(string(jsonBytes))
	if err != nil {
		return err
	}
	b.EncryptedBindingParameters = ciphertext
	return nil
}

// GetBindingParameters decrypts the EncryptedBindingParameters field and
// unmarshals the result into the provided bindingParameters object
func (b *Binding) GetBindingParameters(
	params interface{},
	codec crypto.Codec,
) error {
	if b.EncryptedBindingParameters == "" {
		return nil
	}
	plaintext, err := codec.Decrypt(b.EncryptedBindingParameters)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(plaintext), params)
	if err != nil {
		return err
	}
	return nil
}

// SetBindingContext marshals the provided bindingContext object, encrypts the
// result, and stores it in the EncryptedBindingContext field
func (b *Binding) SetBindingContext(
	context interface{},
	codec crypto.Codec,
) error {
	jsonBytes, err := json.Marshal(context)
	if err != nil {
		return err
	}
	ciphertext, err := codec.Encrypt(string(jsonBytes))
	if err != nil {
		return err
	}
	b.EncryptedBindingContext = string(ciphertext)
	return nil
}

// GetBindingContext decrypts the EncryptedBindingContext field and unmarshals
// the result into the provided bindingContext object
func (b *Binding) GetBindingContext(
	context interface{},
	codec crypto.Codec,
) error {
	if b.EncryptedBindingContext == "" {
		return nil
	}
	plaintext, err := codec.Decrypt(b.EncryptedBindingContext)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(plaintext), context)
	if err != nil {
		return err
	}
	return nil
}

// SetCredentials marshals the provided credentials object, encrypts the result,
// and stores it in the EncryptedCredentials field
func (b *Binding) SetCredentials(
	credentials interface{},
	codec crypto.Codec,
) error {
	jsonBytes, err := json.Marshal(credentials)
	if err != nil {
		return err
	}
	ciphertext, err := codec.Encrypt(string(jsonBytes))
	if err != nil {
		return err
	}
	b.EncryptedCredentials = string(ciphertext)
	return nil
}

// GetCredentials decrypts the EncryptedCredentials field and unmarshals the
// result into the provided credentials object
func (b *Binding) GetCredentials(
	credentials interface{},
	codec crypto.Codec,
) error {
	if b.EncryptedCredentials == "" {
		return nil
	}
	plaintext, err := codec.Decrypt(b.EncryptedCredentials)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(plaintext), credentials)
	if err != nil {
		return err
	}
	return nil
}
