package service

import (
	"encoding/json"
	"time"
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
	Created                time.Time               `json:"created"`
}

// NewInstanceFromJSON returns a new Instance unmarshalled from the provided
// JSON []byte
func NewInstanceFromJSON(
	jsonBytes []byte,
	emptyInstanceDetails InstanceDetails,
	provisioningParametersSchema *InputParametersSchema, // nolint: interfacer
) (Instance, error) {
	instance := Instance{
		Details: emptyInstanceDetails,
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
	}
	err := json.Unmarshal(jsonBytes, &instance)
	return instance, err
}

// ToJSON returns a []byte containing a JSON representation of the
// instance
func (i Instance) ToJSON() ([]byte, error) {
	return json.Marshal(i)
}
