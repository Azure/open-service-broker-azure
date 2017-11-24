package api

import (
	"encoding/json"
)

// BindingRequest represents a request to bind to a service
type BindingRequest struct {
	ServiceID  string                 `json:"service_id"`
	PlanID     string                 `json:"plan_id"`
	Parameters map[string]interface{} `json:"parameters"`
}

// NewBindingRequestFromJSON returns a new BindingRequest unmarshaled from the
// provided JSON []byte
func NewBindingRequestFromJSON(
	jsonBytes []byte,
) (*BindingRequest, error) {
	bindingRequest := &BindingRequest{}
	err := json.Unmarshal(jsonBytes, bindingRequest)
	if err != nil {
		return nil, err
	}
	return bindingRequest, nil
}

// ToJSON returns a []byte containing a JSON representation of the binding
// request
func (b *BindingRequest) ToJSON() ([]byte, error) {
	return json.Marshal(b)
}
