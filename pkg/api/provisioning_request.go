package api

import (
	"encoding/json"

	"github.com/Azure/azure-service-broker/pkg/service"
)

// ProvisioningRequest represents a request to provision a service
type ProvisioningRequest struct {
	ServiceID  string                         `json:"service_id"`
	PlanID     string                         `json:"plan_id"`
	Parameters service.ProvisioningParameters `json:"parameters"`
}

// GetProvisioningRequestFromJSON populates the given ProvisioningRequest by
// unmarshalling the provided JSON []byte
func GetProvisioningRequestFromJSON(
	jsonBytes []byte,
	provisioningRequest *ProvisioningRequest,
) error {
	return json.Unmarshal(jsonBytes, provisioningRequest)
}

// ToJSON returns a []byte containing a JSON representation of the provisioning
// request
func (p *ProvisioningRequest) ToJSON() ([]byte, error) {
	return json.Marshal(p)
}
