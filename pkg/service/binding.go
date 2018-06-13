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
	schema *InputParametersSchema, // nolint: interfacer
) (Binding, error) {
	binding := Binding{
		BindingParameters: &BindingParameters{
			Parameters: Parameters{
				Schema: schema,
			},
		},
		SecureDetails: SecureBindingDetails{},
	}
	if err := json.Unmarshal(jsonBytes, &binding); err != nil {
		return binding, err
	}
	return binding.decrypt()
}

// ToJSON returns a []byte containing a JSON representation of the instance
func (b Binding) ToJSON() ([]byte, error) {
	var err error
	if b, err = b.encrypt(); err != nil {
		return nil, err
	}
	return json.Marshal(b)
}

func (b Binding) encrypt() (Binding, error) {
	return b.encryptSecureDetails()
}

func (b Binding) encryptSecureDetails() (Binding, error) {
	jsonBytes, err := json.Marshal(b.SecureDetails)
	if err != nil {
		return b, err
	}
	b.EncryptedSecureDetails, err = crypto.Encrypt(jsonBytes)
	return b, err
}

func (b Binding) decrypt() (Binding, error) {
	return b.decryptSecureDetails()
}

func (b Binding) decryptSecureDetails() (Binding, error) {
	if len(b.EncryptedSecureDetails) == 0 ||
		b.SecureDetails == nil {
		return b, nil
	}
	plaintext, err := crypto.Decrypt(b.EncryptedSecureDetails)
	if err != nil {
		return b, err
	}
	return b, json.Unmarshal(plaintext, &b.SecureDetails)
}
