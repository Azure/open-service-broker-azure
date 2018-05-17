package api

import (
	"encoding/json"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

// ProvisioningRequest represents a request to provision a service
type ProvisioningRequest struct {
	ServiceID  string             `json:"service_id"`
	PlanID     string             `json:"plan_id"`
	Parameters service.Parameters `json:"parameters"`
}

// NewProvisioningRequestFromJSON returns a new ProvisioningRequest unmarshaled
// from the provided JSON []byte
func NewProvisioningRequestFromJSON(
	jsonBytes []byte,
) (*ProvisioningRequest, error) {
	provisioningRequest := &ProvisioningRequest{}
	err := json.Unmarshal(jsonBytes, provisioningRequest)
	if err != nil {
		return nil, err
	}
	return provisioningRequest, nil
}

// ToJSON returns a []byte containing a JSON representation of the provisioning
// request
func (p *ProvisioningRequest) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}
