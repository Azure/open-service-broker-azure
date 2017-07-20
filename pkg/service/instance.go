package service

import "encoding/json"

// Instance represents an instance of a service
type Instance struct {
	InstanceID string `json:"instanceId"`
	ServiceID  string `json:"serviceId"`
	PlanID     string `json:"planId"`
	// TODO: These should be encrypted as well
	EncodedProvisioningParameters string `json:"provisioningParameters"`
	Status                        string `json:"status"`
	StatusReason                  string `json:"statusReason"`
	// TODO: These should be encrypted as well
	EncodedProvisioningResult string `json:"provisioningResult"`
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

// SetProvisioningParameters marshals the provided parameters and stores them in
// the EncodedProvisioningParameters field
func (i *Instance) SetProvisioningParameters(params interface{}) error {
	bytes, err := json.Marshal(params)
	if err != nil {
		return err
	}
	i.EncodedProvisioningParameters = string(bytes)
	return nil
}

// GetProvisioningParameters unmarshals the EncodedProvisioningParameters into
// the provided object
func (i *Instance) GetProvisioningParameters(params interface{}) error {
	if i.EncodedProvisioningParameters == "" {
		return nil
	}
	err := json.Unmarshal([]byte(i.EncodedProvisioningParameters), params)
	if err != nil {
		return err
	}
	return nil
}

// SetProvisioningResult marshals the provided provisioning result and stores it
// in the EncodedProvisioningResult field
func (i *Instance) SetProvisioningResult(result interface{}) error {
	bytes, err := json.Marshal(result)
	if err != nil {
		return err
	}
	i.EncodedProvisioningResult = string(bytes)
	return nil
}

// GetProvisioningResult unmarshals the EncodedProvisioningResult into the
// provided object
func (i *Instance) GetProvisioningResult(result interface{}) error {
	if i.EncodedProvisioningResult == "" {
		return nil
	}
	err := json.Unmarshal([]byte(i.EncodedProvisioningResult), result)
	if err != nil {
		return err
	}
	return nil
}

// CREATE TABLE instances (
// 	azureInstanceId varchar(256) NOT NULL UNIQUE,
// 	status varchar(18),
// 	timestamp DATETIME DEFAULT (GETDATE()),
// 	instanceId char(36) PRIMARY KEY,
// 	serviceId char(36) NOT NULL,
// 	planId char(36) NOT NULL,
// 	organizationGuid char(36) NOT NULL,
// 	spaceGuid char(36) NOT NULL,
// 	parameters text,
// 	lastOperation text,
// 	provisioningResult text
// );
