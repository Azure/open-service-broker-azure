package service

import "encoding/json"

// Binding represents a biding to a service
type Binding struct {
	BindingID  string `json:"bindingId"`
	InstanceID string `json:"instanceId"`
	// TODO: These should be encrypted as well
	EncodedBindingParameters string `json:"bindingParameters"`
	Status                   string `json:"status"`
	// TODO: These should be encrypted as well
	EncodedBindingResult string `json:"bindingResult"`
}

// NewBindingFromJSONString returns a new Binding unmarshalled from the
// provided JSON string
func NewBindingFromJSONString(jsonStr string) (*Binding, error) {
	binding := &Binding{}
	err := json.Unmarshal([]byte(jsonStr), binding)
	if err != nil {
		return nil, err
	}
	return binding, nil
}

// ToJSONString returns a string containing a JSON representation of the
// instance
func (b *Binding) ToJSONString() (string, error) {
	bytes, err := json.Marshal(b)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// SetBindingParameters marshals the provided parameters and stores them in the
// EncodedBindingParameters field
func (b *Binding) SetBindingParameters(params interface{}) error {
	bytes, err := json.Marshal(params)
	if err != nil {
		return err
	}
	b.EncodedBindingParameters = string(bytes)
	return nil
}

// GetBindingParameters unmarshals the EncodedBindingParameters into the
// provided object
func (b *Binding) GetBindingParameters(params interface{}) error {
	err := json.Unmarshal([]byte(b.EncodedBindingParameters), params)
	if err != nil {
		return err
	}
	return nil
}

// SetBindingResult marshals the provided binding result and stores it in the
// EncodedBindingResult field
func (b *Binding) SetBindingResult(result interface{}) error {
	bytes, err := json.Marshal(result)
	if err != nil {
		return err
	}
	b.EncodedBindingResult = string(bytes)
	return nil
}

// GetBindingResult unmarshals the EncodedBindingResult into the provided object
func (b *Binding) GetBindingResult(result interface{}) error {
	err := json.Unmarshal([]byte(b.EncodedBindingResult), result)
	if err != nil {
		return err
	}
	return nil
}

// CREATE TABLE bindings (
// 	bindingId char(36) PRIMARY KEY,
// 	instanceId char(36) FOREIGN KEY REFERENCES instances(instanceId),
// 	timestamp DATETIME DEFAULT (GETDATE()),
// 	serviceId char(36) NOT NULL,
// 	planId char(36) NOT NULL,
// 	parameters text,
// 	bindingResult text
// );
