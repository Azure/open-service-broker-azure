package service

import (
	"encoding/json"

	"github.com/Azure/azure-service-broker/pkg/crypto"
)

// Instance represents an instance of a service
type Instance struct {
	InstanceID                      string `json:"instanceId"`
	ServiceID                       string `json:"serviceId"`
	PlanID                          string `json:"planId"`
	EncryptedProvisioningParameters string `json:"provisioningParameters"`
	Status                          string `json:"status"`
	StatusReason                    string `json:"statusReason"`
	EncryptedProvisioningContext    string `json:"provisioningContext"`
}

// NewInstanceFromJSONString returns a new Instance unmarshalled from the
// provided JSON string
func NewInstanceFromJSONString(jsonStr string) (*Instance, error) {
	instance := &Instance{}
	err := json.Unmarshal([]byte(jsonStr), instance)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

// ToJSONString returns a string containing a JSON representation of the
// instance
func (i *Instance) ToJSONString() (string, error) {
	bytes, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// SetProvisioningParameters marshals the provided provisioningParameters
// object, encrypts the result, and stores it in the
// EncryptedProvisioningParameters field
func (i *Instance) SetProvisioningParameters(
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
	i.EncryptedProvisioningParameters = ciphertext
	return nil
}

// GetProvisioningParameters decrypts the EncryptedProvisioningParameters field
// and unmarshals the result into the provided provisioningParameters object
func (i *Instance) GetProvisioningParameters(
	params interface{},
	codec crypto.Codec,
) error {
	if i.EncryptedProvisioningParameters == "" {
		return nil
	}
	plaintext, err := codec.Decrypt(i.EncryptedProvisioningParameters)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(plaintext), params)
	if err != nil {
		return err
	}
	return nil
}

// SetProvisioningContext marshals the provided provisioningContext object,
// encrypts the result, and stores it in the EncrypredProvisioningContext field
func (i *Instance) SetProvisioningContext(
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
	i.EncryptedProvisioningContext = ciphertext
	return nil
}

// GetProvisioningContext decrypts the EncryptedProvisioningContext field and
// unmarshals the result into the provided provisioningContext object
func (i *Instance) GetProvisioningContext(
	context interface{},
	codec crypto.Codec,
) error {
	if i.EncryptedProvisioningContext == "" {
		return nil
	}
	plaintext, err := codec.Decrypt(i.EncryptedProvisioningContext)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(plaintext), context)
	if err != nil {
		return err
	}
	return nil
}
