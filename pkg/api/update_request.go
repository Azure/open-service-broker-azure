package api

import (
	"encoding/json"

	"github.com/Azure/azure-service-broker/pkg/service"
)

// UpdatingPreviousValues represents the information about the service instance
// prior to the update. Our broker doesn't need it. Per spec, it still could be
// provided.
type UpdatingPreviousValues struct {
	PlanID string `json:"plan_id"`
}

// UpdatingRequest represents a request to update a service
type UpdatingRequest struct {
	ServiceID      string                     `json:"service_id"`
	PlanID         string                     `json:"plan_id"`
	Parameters     service.UpdatingParameters `json:"parameters"`
	PreviousValues UpdatingPreviousValues     `json:"previous_values"`
}

// GetUpdatingRequestFromJSON populates the given UpdatingRequest by
// unmarshalling the provided JSON []byte
func GetUpdatingRequestFromJSON(
	jsonBytes []byte,
	updatingRequest *UpdatingRequest,
) error {
	return json.Unmarshal(jsonBytes, updatingRequest)
}

// ToJSON returns a []byte containing a JSON representation of the updating
// request
func (u *UpdatingRequest) ToJSON() ([]byte, error) {
	return json.Marshal(u)
}
