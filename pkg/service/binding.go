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
	EncodedBindingContext string `json:"bindingContext"`
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

// SetBindingParameters marshals the providedParameters and stores them in the
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
	if b.EncodedBindingParameters == "" {
		return nil
	}
	err := json.Unmarshal([]byte(b.EncodedBindingParameters), params)
	if err != nil {
		return err
	}
	return nil
}

// SetBindingContext marshals the provided bindingContext and stores it in the
// EncodedBindingContext field
func (b *Binding) SetBindingContext(context interface{}) error {
	bytes, err := json.Marshal(context)
	if err != nil {
		return err
	}
	b.EncodedBindingContext = string(bytes)
	return nil
}

// GetBindingContext unmarshals the EncodedBindingContext into the provided
// object
func (b *Binding) GetBindingContext(context interface{}) error {
	if b.EncodedBindingContext == "" {
		return nil
	}
	err := json.Unmarshal([]byte(b.EncodedBindingContext), context)
	if err != nil {
		return err
	}
	return nil
}
