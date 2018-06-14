package service

import (
	"encoding/json"
	"time"
)

// Binding represents a binding to a service
type Binding struct {
	BindingID         string             `json:"bindingId"`
	InstanceID        string             `json:"instanceId"`
	ServiceID         string             `json:"serviceId"`
	BindingParameters *BindingParameters `json:"bindingParameters"`
	Status            string             `json:"status"`
	StatusReason      string             `json:"statusReason"`
	Details           BindingDetails     `json:"details"`
	Created           time.Time          `json:"created"`
}

// NewBindingFromJSON returns a new Binding unmarshalled from the provided JSON
// []byte
func NewBindingFromJSON(
	jsonBytes []byte,
	emptyBindingDetails BindingDetails,
	schema *InputParametersSchema, // nolint: interfacer
) (Binding, error) {
	binding := Binding{
		Details: emptyBindingDetails,
		BindingParameters: &BindingParameters{
			Parameters: Parameters{
				Schema: schema,
			},
		},
	}
	err := json.Unmarshal(jsonBytes, &binding)
	return binding, err
}

// ToJSON returns a []byte containing a JSON representation of the instance
func (b Binding) ToJSON() ([]byte, error) {
	return json.Marshal(b)
}
