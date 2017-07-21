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
	EncodedProvisioningContext string `json:"provisioningContext"`
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

// SetProvisioningParameters marshals the providedParameters and stores them in
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

// SetProvisioningContext marshals the provided object and stores it in the
// EncodedProvisioningContext field
func (i *Instance) SetProvisioningContext(context interface{}) error {
	bytes, err := json.Marshal(context)
	if err != nil {
		return err
	}
	i.EncodedProvisioningContext = string(bytes)
	return nil
}

// GetProvisioningContext unmarshals the EncodedProvisioningContext into the
// provided object
func (i *Instance) GetProvisioningContext(context interface{}) error {
	if i.EncodedProvisioningContext == "" {
		return nil
	}
	err := json.Unmarshal([]byte(i.EncodedProvisioningContext), context)
	if err != nil {
		return err
	}
	return nil
}
