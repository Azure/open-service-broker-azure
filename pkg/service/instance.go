package service

import (
	"encoding/json"
	"time"

	"github.com/Azure/azure-service-broker/pkg/crypto"
)

// Instance represents an instance of a service
type Instance struct {
	InstanceID                      string    `json:"instanceId"`
	ServiceID                       string    `json:"serviceId"`
	PlanID                          string    `json:"planId"`
	EncryptedProvisioningParameters []byte    `json:"provisioningParameters"`
	Status                          string    `json:"status"`
	StatusReason                    string    `json:"statusReason"`
	EncryptedProvisioningContext    []byte    `json:"provisioningContext"`
	Created                         time.Time `json:"created"`
}

// NewInstanceFromJSON returns a new Instance unmarshalled from the provided
// JSON []byte
func NewInstanceFromJSON(jsonBytes []byte) (*Instance, error) {
	instance := &Instance{}
	if err := json.Unmarshal(jsonBytes, instance); err != nil {
		return nil, err
	}
	return instance, nil
}

// ToJSON returns a []byte containing a JSON representation of the
// instance
func (i *Instance) ToJSON() ([]byte, error) {
	return json.Marshal(i)
}

// SetProvisioningParameters marshals the provided provisioningParameters
// object, encrypts the result, and stores it in the
// EncryptedProvisioningParameters field
func (i *Instance) SetProvisioningParameters(
	params ProvisioningParameters,
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
	i.EncryptedProvisioningParameters = ciphertext
	return nil
}

// GetProvisioningParameters decrypts the EncryptedProvisioningParameters field
// and unmarshals the result into the provided provisioningParameters object
func (i *Instance) GetProvisioningParameters(
	params ProvisioningParameters,
	codec crypto.Codec,
) error {
	if len(i.EncryptedProvisioningParameters) == 0 {
		return nil
	}
	plaintext, err := codec.Decrypt(i.EncryptedProvisioningParameters)
	if err != nil {
		return err
	}
	return json.Unmarshal(plaintext, params)
}

// SetProvisioningContext marshals the provided provisioningContext object,
// encrypts the result, and stores it in the EncrypredProvisioningContext field
func (i *Instance) SetProvisioningContext(
	context ProvisioningContext,
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
	i.EncryptedProvisioningContext = ciphertext
	return nil
}

// GetProvisioningContext decrypts the EncryptedProvisioningContext field and
// unmarshals the result into the provided provisioningContext object
func (i *Instance) GetProvisioningContext(
	context ProvisioningContext,
	codec crypto.Codec,
) error {
	if len(i.EncryptedProvisioningContext) == 0 {
		return nil
	}
	plaintext, err := codec.Decrypt(i.EncryptedProvisioningContext)
	if err != nil {
		return err
	}
	return json.Unmarshal(plaintext, context)
}
