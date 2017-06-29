package service

import "encoding/json"

// BindingRequest represents a request to bind to a service
type BindingRequest struct {
	ServiceID  string      `json:"service_id"`
	PlanID     string      `json:"plan_id"`
	Parameters interface{} `json:"parameters"`
}

// GetBindingRequestFromJSONString populates the given BindingRequest by
// unmarshalling the provided JSON string
func GetBindingRequestFromJSONString(
	jsonStr string,
	bindingRequest *BindingRequest,
) error {
	err := json.Unmarshal([]byte(jsonStr), bindingRequest)
	if err != nil {
		return err
	}
	return nil
}
