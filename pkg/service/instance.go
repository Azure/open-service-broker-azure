package service

import (
	"encoding/json"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/crypto"
)

// Instance represents an instance of a service
type Instance struct {
	InstanceID             string                  `json:"instanceId"`
	Alias                  string                  `json:"alias"`
	ServiceID              string                  `json:"serviceId"`
	Service                Service                 `json:"-"`
	PlanID                 string                  `json:"planId"`
	Plan                   Plan                    `json:"-"`
	ProvisioningParameters *ProvisioningParameters `json:"provisioningParameters"`
	UpdatingParameters     *ProvisioningParameters `json:"updatingParameters"`
	Status                 string                  `json:"status"`
	StatusReason           string                  `json:"statusReason"`
	Parent                 *Instance               `json:"-"`
	ParentAlias            string                  `json:"parentAlias"`
	Details                InstanceDetails         `json:"details"`
	EncryptedSecureDetails []byte                  `json:"secureDetails"`
	SecureDetails          SecureInstanceDetails   `json:"-"`
	Created                time.Time               `json:"created"`
}

// NewInstanceFromJSON returns a new Instance unmarshalled from the provided
// JSON []byte
func NewInstanceFromJSON(
	jsonBytes []byte,
	provisioningParametersSchema *InputParametersSchema, // nolint: interfacer
) (Instance, error) {
	instance := Instance{
		ProvisioningParameters: &ProvisioningParameters{
			Parameters: Parameters{
				Schema: provisioningParametersSchema,
			},
		},
		UpdatingParameters: &ProvisioningParameters{
			Parameters: Parameters{
				// Note that provisioning schema is deliberately used here in place of
				// updating schema. That allows us to store/retrieve the FULL set of
				// combined provisioning + updating parameters and not just the subset
				// of provisioning parameters that are also valid updating parameters.
				Schema: provisioningParametersSchema,
			},
		},
		Details:       InstanceDetails{},
		SecureDetails: SecureInstanceDetails{},
	}
	if err := json.Unmarshal(jsonBytes, &instance); err != nil {
		return instance, err
	}
	return instance.decrypt()
}

// ToJSON returns a []byte containing a JSON representation of the
// instance
func (i Instance) ToJSON() ([]byte, error) {
	var err error
	if i, err = i.encrypt(); err != nil {
		return nil, err
	}
	return json.Marshal(i)
}

func (i Instance) encrypt() (Instance, error) {
	return i.encryptSecureDetails()
}

func (i Instance) encryptSecureDetails() (Instance, error) {
	jsonBytes, err := json.Marshal(i.SecureDetails)
	if err != nil {
		return i, err
	}
	i.EncryptedSecureDetails, err = crypto.Encrypt(jsonBytes)
	return i, err
}

func (i Instance) decrypt() (Instance, error) {
	return i.decryptSecureDetails()
}

func (i Instance) decryptSecureDetails() (Instance, error) {
	if len(i.EncryptedSecureDetails) == 0 ||
		i.SecureDetails == nil {
		return i, nil
	}
	plaintext, err := crypto.Decrypt(i.EncryptedSecureDetails)
	if err != nil {
		return i, err
	}
	return i, json.Unmarshal(plaintext, &i.SecureDetails)
}
