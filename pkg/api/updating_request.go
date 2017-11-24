package api

import (
	"encoding/json"
)

// UpdatingPreviousValues represents the information about the service instance
// prior to the update. Our broker doesn't need it. Per spec, it still could be
// provided.
type UpdatingPreviousValues struct {
	PlanID string `json:"plan_id"`
}

// UpdatingRequest represents a request to update a service
type UpdatingRequest struct {
	ServiceID      string                 `json:"service_id"`
	PlanID         string                 `json:"plan_id"`
	Parameters     map[string]interface{} `json:"parameters"`
	PreviousValues UpdatingPreviousValues `json:"previous_values"`
}

// NewUpdatingRequestFromJSON returns a new UpdatingRequest unmarshaled from the
// provided JSON []byte
func NewUpdatingRequestFromJSON(
	jsonBytes []byte,
) (*UpdatingRequest, error) {
	updatingRequest := &UpdatingRequest{}
	err := json.Unmarshal(jsonBytes, updatingRequest)
	if err != nil {
		return nil, err
	}
	return updatingRequest, nil
}

// ToJSON returns a []byte containing a JSON representation of the updating
// request
func (u *UpdatingRequest) ToJSON() ([]byte, error) {
	return json.Marshal(u)
}
