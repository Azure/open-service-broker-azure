package service

import (
	"encoding/json"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/crypto"
)

// Instance represents an instance of a service
type Instance struct {
	InstanceID                      string                 `json:"instanceId"`
	Alias                           string                 `json:"alias"`
	ServiceID                       string                 `json:"serviceId"`
	Service                         Service                `json:"-"`
	PlanID                          string                 `json:"planId"`
	Plan                            Plan                   `json:"-"`
	EncryptedProvisioningParameters []byte                 `json:"provisioningParameters"` // nolint: lll
	ProvisioningParameters          ProvisioningParameters `json:"-"`
	EncryptedUpdatingParameters     []byte                 `json:"updatingParameters"` // nolint: lll
	UpdatingParameters              UpdatingParameters     `json:"-"`
	Status                          string                 `json:"status"`
	StatusReason                    string                 `json:"statusReason"`
	Location                        string                 `json:"location"`
	ResourceGroup                   string                 `json:"resourceGroup"`
	Parent                          *Instance              `json:"-"`
	ParentAlias                     string                 `json:"parentAlias"`
	Tags                            map[string]string      `json:"tags"`
	Details                         InstanceDetails        `json:"details"`
	EncryptedSecureDetails          []byte                 `json:"secureDetails"`
	SecureDetails                   SecureInstanceDetails  `json:"-"`
	Created                         time.Time              `json:"created"`
}

// NewInstanceFromJSON returns a new Instance unmarshalled from the provided
// JSON []byte
func NewInstanceFromJSON(
	jsonBytes []byte,
	pp ProvisioningParameters,
	up UpdatingParameters,
	dt InstanceDetails,
	sdt InstanceDetails,
	codec crypto.Codec,
) (Instance, error) {
	instance := Instance{
		ProvisioningParameters: pp,
		UpdatingParameters:     up,
		Details:                dt,
		SecureDetails:          sdt,
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
	return i.encryptSecureDetails(codec)
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

func (i Instance) encryptSecureDetails(
	codec crypto.Codec,
) (Instance, error) {
	jsonBytes, err := json.Marshal(i.SecureDetails)
	if err != nil {
		return i, err
	}
	i.EncryptedSecureDetails, err = codec.Encrypt(jsonBytes)
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
	return i.decryptSecureDetails(codec)
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

func (i Instance) decryptSecureDetails(
	codec crypto.Codec,
) (Instance, error) {
	if len(i.EncryptedSecureDetails) == 0 ||
		i.SecureDetails == nil {
		return i, nil
	}
	plaintext, err := codec.Decrypt(i.EncryptedSecureDetails)
	if err != nil {
		return i, err
	}
	return i, json.Unmarshal(plaintext, i.SecureDetails)
}
