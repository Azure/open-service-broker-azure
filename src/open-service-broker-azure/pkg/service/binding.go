package service

import (
	"encoding/json"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/crypto"
)

// Binding represents a binding to a service
type Binding struct {
	BindingID              string               `json:"bindingId"`
	InstanceID             string               `json:"instanceId"`
	ServiceID              string               `json:"serviceId"`
	BindingParameters      *BindingParameters   `json:"bindingParameters"`
	Status                 string               `json:"status"`
	StatusReason           string               `json:"statusReason"`
	Details                BindingDetails       `json:"details"`
	EncryptedSecureDetails []byte               `json:"secureDetails"`
	SecureDetails          SecureBindingDetails `json:"-"`
	Created                time.Time            `json:"created"`
}

// NewBindingFromJSON returns a new Binding unmarshalled from the provided JSON
// []byte
func NewBindingFromJSON(
	jsonBytes []byte,
	codec crypto.Codec,
	schema *InputParametersSchema, // nolint: interfacer
) (Binding, error) {
	binding := Binding{
		BindingParameters: &BindingParameters{
			Parameters: Parameters{
				Codec:  codec,
				Schema: schema,
			},
		},
		SecureDetails: SecureBindingDetails{},
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
	// Set the codec on the params before continuing
	if b.BindingParameters != nil {
		b.BindingParameters.Codec = codec
	}
	return json.Marshal(b)
}

func (b Binding) encrypt(codec crypto.Codec) (Binding, error) {
	return b.encryptSecureDetails(codec)
}

func (b Binding) encryptSecureDetails(codec crypto.Codec) (Binding, error) {
	jsonBytes, err := json.Marshal(b.SecureDetails)
	if err != nil {
		return b, err
	}
	b.EncryptedSecureDetails, err = codec.Encrypt(jsonBytes)
	return b, err
}

func (b Binding) decrypt(codec crypto.Codec) (Binding, error) {
	return b.decryptSecureDetails(codec)
}

func (b Binding) decryptSecureDetails(codec crypto.Codec) (Binding, error) {
	if len(b.EncryptedSecureDetails) == 0 ||
		b.SecureDetails == nil {
		return b, nil
	}
	plaintext, err := codec.Decrypt(b.EncryptedSecureDetails)
	if err != nil {
		return b, err
	}
	return b, json.Unmarshal(plaintext, &b.SecureDetails)
}
