package service

import (
	"encoding/json"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/crypto"
)

// Binding represents a binding to a service
type Binding struct {
	BindingID                  string            `json:"bindingId"`
	InstanceID                 string            `json:"instanceId"`
	ServiceID                  string            `json:"serviceId"`
	EncryptedBindingParameters []byte            `json:"bindingParameters"`
	BindingParameters          BindingParameters `json:"-"`
	Status                     string            `json:"status"`
	StatusReason               string            `json:"statusReason"`
	EncryptedBindingContext    []byte            `json:"bindingContext"`
	BindingContext             BindingContext    `json:"-"`
	EncryptedCredentials       []byte            `json:"credentials"`
	Credentials                Credentials       `json:"-"`
	Created                    time.Time         `json:"created"`
}

// NewBindingFromJSON returns a new Binding unmarshalled from the provided JSON
// []byte
func NewBindingFromJSON(
	jsonBytes []byte,
	bp BindingParameters,
	bc BindingContext,
	cr Credentials,
	codec crypto.Codec,
) (Binding, error) {
	binding := Binding{
		BindingParameters: bp,
		BindingContext:    bc,
		Credentials:       cr,
	}
	if err := json.Unmarshal(jsonBytes, &binding); err != nil {
		return binding, err
	}
	return binding.decrypt(codec)
}

// ToJSON returns a []byte containing a JSON representation of the instance
func (b Binding) ToJSON(codec crypto.Codec) ([]byte, error) {
	var err error
	if b, err = b.encrypt(codec); err != nil {
		return nil, err
	}
	return json.Marshal(b)
}

func (b Binding) encrypt(codec crypto.Codec) (Binding, error) {
	var err error
	if b, err = b.encryptBindingParameters(codec); err != nil {
		return b, err
	}
	if b, err = b.encryptBindingContext(codec); err != nil {
		return b, err
	}
	return b.encryptCredentials(codec)
}

func (b Binding) encryptBindingParameters(codec crypto.Codec) (Binding, error) {
	jsonBytes, err := json.Marshal(b.BindingParameters)
	if err != nil {
		return b, err
	}
	b.EncryptedBindingParameters, err = codec.Encrypt(jsonBytes)
	return b, err
}

func (b Binding) encryptBindingContext(codec crypto.Codec) (Binding, error) {
	jsonBytes, err := json.Marshal(b.BindingContext)
	if err != nil {
		return b, err
	}
	b.EncryptedBindingContext, err = codec.Encrypt(jsonBytes)
	return b, err
}

func (b Binding) encryptCredentials(codec crypto.Codec) (Binding, error) {
	jsonBytes, err := json.Marshal(b.Credentials)
	if err != nil {
		return b, err
	}
	b.EncryptedCredentials, err = codec.Encrypt(jsonBytes)
	return b, err
}

func (b Binding) decrypt(codec crypto.Codec) (Binding, error) {
	var err error
	if b, err = b.decryptBindingParameters(codec); err != nil {
		return b, err
	}
	if b, err = b.decryptBindingContext(codec); err != nil {
		return b, err
	}
	return b.decryptCredentials(codec)
}

func (b Binding) decryptBindingParameters(codec crypto.Codec) (Binding, error) {
	if len(b.EncryptedBindingParameters) == 0 ||
		b.BindingParameters == nil {
		return b, nil
	}
	plaintext, err := codec.Decrypt(b.EncryptedBindingParameters)
	if err != nil {
		return b, err
	}
	return b, json.Unmarshal(plaintext, b.BindingParameters)
}

func (b Binding) decryptBindingContext(codec crypto.Codec) (Binding, error) {
	if len(b.EncryptedBindingContext) == 0 ||
		b.BindingContext == nil {
		return b, nil
	}
	plaintext, err := codec.Decrypt(b.EncryptedBindingContext)
	if err != nil {
		return b, err
	}
	return b, json.Unmarshal(plaintext, b.BindingContext)
}

func (b Binding) decryptCredentials(codec crypto.Codec) (Binding, error) {
	if len(b.EncryptedCredentials) == 0 ||
		b.Credentials == nil {
		return b, nil
	}
	plaintext, err := codec.Decrypt(b.EncryptedCredentials)
	if err != nil {
		return b, err
	}
	return b, json.Unmarshal(plaintext, b.Credentials)
}
