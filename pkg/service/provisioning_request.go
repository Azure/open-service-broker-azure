package service

import "encoding/json"

// ProvisioningRequest represents a request to provision a service
type ProvisioningRequest struct {
	ServiceID  string      `json:"service_id"`
	PlanID     string      `json:"plan_id"`
	Parameters interface{} `json:"parameters"`
}

// GetProvisioningRequestFromJSONString populates the given ProvisioningRequest
// by unmarshalling the provided JSON string
func GetProvisioningRequestFromJSONString(
	jsonStr string,
	provisioningRequest *ProvisioningRequest,
) error {
	err := json.Unmarshal([]byte(jsonStr), provisioningRequest)
	if err != nil {
		return err
	}
	return nil
}
