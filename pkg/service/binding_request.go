package service

import "encoding/json"

// BindingRequest represents a request to bind to a service
type BindingRequest struct {
	ServiceID  string      `json:"service_id"`
	PlanID     string      `json:"plan_id"`
	Parameters interface{} `json:"parameters"`
}

// GetBindingRequestFromJSON populates the given BindingRequest by unmarshalling
// the provided JSON []byte
func GetBindingRequestFromJSON(
	jsonBytes []byte,
	bindingRequest *BindingRequest,
) error {
	return json.Unmarshal(jsonBytes, bindingRequest)
}

// ToJSON returns a []byte containing a JSON representation of the binding
// request
func (b *BindingRequest) ToJSON() ([]byte, error) {
	return json.Marshal(b)
}
