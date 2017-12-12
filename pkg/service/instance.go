package service

import (
	"encoding/json"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/crypto"
)

// Instance represents an instance of a service
type Instance struct {
	InstanceID                      string                         `json:"instanceId"`                     // nolint: lll
	ServiceID                       string                         `json:"serviceId"`                      // nolint: lll
	PlanID                          string                         `json:"planId"`                         // nolint: lll
	StandardProvisioningParameters  StandardProvisioningParameters `json:"standardProvisioningParameters"` // nolint: lll
	EncryptedProvisioningParameters []byte                         `json:"provisioningParameters"`         // nolint: lll
	EncryptedUpdatingParameters     []byte                         `json:"updatingParameters"`             // nolint: lll
	Status                          string                         `json:"status"`                         // nolint: lll
	StatusReason                    string                         `json:"statusReason"`                   // nolint: lll
	StandardProvisioningContext     StandardProvisioningContext    `json:"standardProvisioningContext"`    // nolint: lll
	EncryptedProvisioningContext    []byte                         `json:"provisioningContext"`            // nolint: lll
	Created                         time.Time                      `json:"created"`                        // nolint: lll
}

// NewInstanceFromJSON returns a new Instance unmarshalled from the provided
// JSON []byte
func NewInstanceFromJSON(jsonBytes []byte) (Instance, error) {
	instance := Instance{}
	err := json.Unmarshal(jsonBytes, &instance)
	return instance, err
}

// ToJSON returns a []byte containing a JSON representation of the
// instance
func (i Instance) ToJSON() ([]byte, error) {
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

// SetUpdatingParameters marshals the provided updatingParameters
// object, encrypts the result, and stores it in the
// EncryptedUpdatingParameters field
func (i *Instance) SetUpdatingParameters(
	params UpdatingParameters,
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
	i.EncryptedUpdatingParameters = ciphertext
	return nil
}

// GetProvisioningParameters decrypts the EncryptedProvisioningParameters field
// and unmarshals the result into the provided provisioningParameters object
func (i Instance) GetProvisioningParameters(
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

// GetUpdatingParameters decrypts the EncryptedUpdatingParameters field
// and unmarshals the result into the provided updatingParameters object
func (i Instance) GetUpdatingParameters(
	params UpdatingParameters,
	codec crypto.Codec,
) error {
	if len(i.EncryptedUpdatingParameters) == 0 {
		return nil
	}
	plaintext, err := codec.Decrypt(i.EncryptedUpdatingParameters)
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
func (i Instance) GetProvisioningContext(
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
