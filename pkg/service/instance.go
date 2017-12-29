package service

import (
	"encoding/json"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/crypto"
)

// Instance represents an instance of a service
type Instance struct {
	InstanceID                      string                 `json:"instanceId"` // nolint: lll
	ServiceID                       string                 `json:"serviceId"`  // nolint: lll
	PlanID                          string                 `json:"planId"`     // nolint: lll
	Parent                          *Instance              `json:"-"`
	EncryptedProvisioningParameters []byte                 `json:"provisioningParameters"` // nolint: lll
	ProvisioningParameters          ProvisioningParameters `json:"-"`
	EncryptedUpdatingParameters     []byte                 `json:"updatingParameters"` // nolint: lll
	UpdatingParameters              UpdatingParameters     `json:"-"`
	Status                          string                 `json:"status"`        // nolint: lll
	StatusReason                    string                 `json:"statusReason"`  // nolint: lll
	Location                        string                 `json:"location"`      // nolint: lll
	ResourceGroup                   string                 `json:"resourceGroup"` // nolint: lll
	Tags                            map[string]string      `json:"tags"`
	EncryptedDetails                []byte                 `json:"details"` // nolint: lll
	Details                         InstanceDetails        `json:"-"`
	Created                         time.Time              `json:"created"` // nolint: lll
}

// NewInstanceFromJSON returns a new Instance unmarshalled from the provided
// JSON []byte
func NewInstanceFromJSON(
	jsonBytes []byte,
	pp ProvisioningParameters,
	up UpdatingParameters,
	dt InstanceDetails,
	codec crypto.Codec,
) (Instance, error) {
	instance := Instance{
		ProvisioningParameters: pp,
		UpdatingParameters:     up,
		Details:                dt,
	}
	if err := json.Unmarshal(jsonBytes, &instance); err != nil {
		return instance, err
	}
	return instance.decrypt(codec)
}

// ToJSON returns a []byte containing a JSON representation of the
// instance
func (i Instance) ToJSON(codec crypto.Codec) ([]byte, error) {
	var err error
	if i, err = i.encrypt(codec); err != nil {
		return nil, err
	}
	return json.Marshal(i)
}

func (i Instance) encrypt(codec crypto.Codec) (Instance, error) {
	var err error
	if i, err = i.encryptProvisioningParameters(codec); err != nil {
		return i, err
	}
	if i, err = i.encryptUpdatingParameters(codec); err != nil {
		return i, err
	}
	return i.encryptDetails(codec)
}

func (i Instance) encryptProvisioningParameters(
	codec crypto.Codec,
) (Instance, error) {
	jsonBytes, err := json.Marshal(i.ProvisioningParameters)
	if err != nil {
		return i, err
	}
	i.EncryptedProvisioningParameters, err = codec.Encrypt(jsonBytes)
	return i, err
}

func (i Instance) encryptUpdatingParameters(
	codec crypto.Codec,
) (Instance, error) {
	jsonBytes, err := json.Marshal(i.UpdatingParameters)
	if err != nil {
		return i, err
	}
	i.EncryptedUpdatingParameters, err = codec.Encrypt(jsonBytes)
	return i, err
}

func (i Instance) encryptDetails(
	codec crypto.Codec,
) (Instance, error) {
	jsonBytes, err := json.Marshal(i.Details)
	if err != nil {
		return i, err
	}
	i.EncryptedDetails, err = codec.Encrypt(jsonBytes)
	return i, err
}

func (i Instance) decrypt(codec crypto.Codec) (Instance, error) {
	var err error
	if i, err = i.decryptProvisioningParameters(codec); err != nil {
		return i, err
	}
	if i, err = i.decryptUpdatingParameters(codec); err != nil {
		return i, err
	}
	return i.decryptDetails(codec)
}

func (i Instance) decryptProvisioningParameters(
	codec crypto.Codec,
) (Instance, error) {
	if len(i.EncryptedProvisioningParameters) == 0 ||
		i.ProvisioningParameters == nil {
		return i, nil
	}
	plaintext, err := codec.Decrypt(i.EncryptedProvisioningParameters)
	if err != nil {
		return i, err
	}
	return i, json.Unmarshal(plaintext, i.ProvisioningParameters)
}

func (i Instance) decryptUpdatingParameters(
	codec crypto.Codec,
) (Instance, error) {
	if len(i.EncryptedUpdatingParameters) == 0 ||
		i.UpdatingParameters == nil {
		return i, nil
	}
	plaintext, err := codec.Decrypt(i.EncryptedUpdatingParameters)
	if err != nil {
		return i, err
	}
	return i, json.Unmarshal(plaintext, i.UpdatingParameters)
}

func (i Instance) decryptDetails(
	codec crypto.Codec,
) (Instance, error) {
	if len(i.EncryptedDetails) == 0 ||
		i.Details == nil {
		return i, nil
	}
	plaintext, err := codec.Decrypt(i.EncryptedDetails)
	if err != nil {
		return i, err
	}
	return i, json.Unmarshal(plaintext, i.Details)
}
