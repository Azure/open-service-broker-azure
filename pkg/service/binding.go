package service

import (
	"encoding/json"

	"github.com/Azure/azure-service-broker/pkg/crypto"
)

// Binding represents a biding to a service
type Binding struct {
	BindingID                  string `json:"bindingId"`
	InstanceID                 string `json:"instanceId"`
	EncryptedBindingParameters []byte `json:"bindingParameters"`
	Status                     string `json:"status"`
	StatusReason               string `json:"statusReason"`
	EncryptedBindingContext    []byte `json:"bindingContext"`
	EncryptedCredentials       []byte `json:"credentials"`
}

// NewBindingFromJSON returns a new Binding unmarshalled from the provided JSON
// []byte
func NewBindingFromJSON(jsonBytes []byte) (*Binding, error) {
	binding := &Binding{}
	if err := json.Unmarshal(jsonBytes, binding); err != nil {
		return nil, err
	}
	return binding, nil
}

// ToJSON returns a []byte containing a JSON representation of the instance
func (b *Binding) ToJSON() ([]byte, error) {
	return json.Marshal(b)
}

// SetBindingParameters marshals the provided bindingParameters object, encrypts
// the result, and stores it in the EncryptedBindingParameters field
func (b *Binding) SetBindingParameters(
	params BindingParameters,
	codec crypto.Codec,
) error {
	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return err
	}
	ciphertext, err := codec.Encrypt(jsonBytes)
	if err != nil {
		return err
	}
	b.EncryptedBindingParameters = ciphertext
	return nil
}

// GetBindingParameters decrypts the EncryptedBindingParameters field and
// unmarshals the result into the provided bindingParameters object
func (b *Binding) GetBindingParameters(
	params BindingParameters,
	codec crypto.Codec,
) error {
	if len(b.EncryptedBindingParameters) == 0 {
		return nil
	}
	plaintext, err := codec.Decrypt(b.EncryptedBindingParameters)
	if err != nil {
		return err
	}
	return json.Unmarshal(plaintext, params)
}

// SetBindingContext marshals the provided bindingContext object, encrypts the
// result, and stores it in the EncryptedBindingContext field
func (b *Binding) SetBindingContext(
	context BindingContext,
	codec crypto.Codec,
) error {
	jsonBytes, err := json.Marshal(context)
	if err != nil {
		return err
	}
	ciphertext, err := codec.Encrypt(jsonBytes)
	if err != nil {
		return err
	}
	b.EncryptedBindingContext = ciphertext
	return nil
}

// GetBindingContext decrypts the EncryptedBindingContext field and unmarshals
// the result into the provided bindingContext object
func (b *Binding) GetBindingContext(
	context BindingContext,
	codec crypto.Codec,
) error {
	if len(b.EncryptedBindingContext) == 0 {
		return nil
	}
	plaintext, err := codec.Decrypt(b.EncryptedBindingContext)
	if err != nil {
		return err
	}
	return json.Unmarshal(plaintext, context)
}

// SetCredentials marshals the provided credentials object, encrypts the result,
// and stores it in the EncryptedCredentials field
func (b *Binding) SetCredentials(
	credentials Credentials,
	codec crypto.Codec,
) error {
	jsonBytes, err := json.Marshal(credentials)
	if err != nil {
		return err
	}
	ciphertext, err := codec.Encrypt(jsonBytes)
	if err != nil {
		return err
	}
	b.EncryptedCredentials = ciphertext
	return nil
}

// GetCredentials decrypts the EncryptedCredentials field and unmarshals the
// result into the provided credentials object
func (b *Binding) GetCredentials(
	credentials Credentials,
	codec crypto.Codec,
) error {
	if len(b.EncryptedCredentials) == 0 {
		return nil
	}
	plaintext, err := codec.Decrypt(b.EncryptedCredentials)
	if err != nil {
		return err
	}
	return json.Unmarshal(plaintext, credentials)
}
