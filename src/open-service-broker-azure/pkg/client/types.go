package client

const (
	// OperationProvisioning represents the "provisioning" operation
	OperationProvisioning = "provisioning"
	// OperationUpdating represents the "updating" operation
	OperationUpdating = "updating"
	// OperationDeprovisioning represents the "deprovisioning" operation
	OperationDeprovisioning = "deprovisioning"
	// OperationStateInProgress represents the state of an operation that is still
	// pending completion
	OperationStateInProgress = "in progress"
	// OperationStateSucceeded represents the state of an operation that has
	// completed successfully
	OperationStateSucceeded = "succeeded"
	// OperationStateFailed represents the state of an operation that has
	// failed
	OperationStateFailed = "failed"
	// OperationStateGone is a pseudo oepration state represting the "state"
	// of an operation against an entity that no longer exists
	OperationStateGone = "gone"
)

// Catalog represents the full set of service/plan offerings provided by a
// broker
type Catalog struct {
	Services []Service `json:"services"`
}

// Service represents a KIND of service
type Service struct {
	Name          string           `json:"name"`
	ID            string           `json:"id"`
	Description   string           `json:"description"`
	Metadata      *ServiceMetadata `json:"metadata,omitempty"`
	Plans         []Plan           `json:"plans"`
	Tags          []string         `json:"tags"`
	Bindable      bool             `json:"bindable"`
	PlanUpdatable bool             `json:"plan_updateable"` // Misspelling is
	// deliberate to match the spec
}

// ServiceMetadata represents optional, extended metadata for a service
type ServiceMetadata struct {
	DisplayName         string `json:"displayName,omitempty"`
	ImageURL            string `json:"imageUrl,omitempty"`
	LongDescription     string `json:"longDescription,omitempty"`
	ProviderDisplayName string `json:"providerDisplayName,omitempty"`
	DocumentationURL    string `json:"documentationUrl,omitempty"`
	SupportURL          string `json:"supportUrl,omitempty"`
}

// Plan represents a concrete variation of a Service; typically these are
// different SKUs for the same KIND of service
type Plan struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Free        bool                 `json:"free"`
	Metadata    *ServicePlanMetadata `json:"metadata,omitempty"`
	// TODO: This currently omits schema information
}

// ServicePlanMetadata represents optional, extended metadata for a plan
type ServicePlanMetadata struct {
	DisplayName string   `json:"displayName,omitempty"`
	Bullets     []string `json:"bullets,omitempty"`
}

// ProvisioningRequest represents a request to provision a service instance
type ProvisioningRequest struct {
	ServiceID  string                 `json:"service_id"`
	PlanID     string                 `json:"plan_id"`
	Parameters map[string]interface{} `json:"parameters"`
}

// UpdatingRequest represents a request to update a service
type UpdatingRequest struct {
	ServiceID      string                 `json:"service_id"`
	PlanID         string                 `json:"plan_id"`
	Parameters     map[string]interface{} `json:"parameters"`
	PreviousValues UpdatingPreviousValues `json:"previous_values"`
}

// UpdatingPreviousValues represents the information about the service instance
// prior to the update
type UpdatingPreviousValues struct {
	PlanID string `json:"plan_id"`
}

// BindingRequest represents a request to bind to a service
type BindingRequest struct {
	ServiceID  string                 `json:"service_id"`
	PlanID     string                 `json:"plan_id"`
	Parameters map[string]interface{} `json:"parameters"`
}

// BindingResponse represents the response to a binding request
type BindingResponse struct {
	Credentials map[string]interface{} `json:"credentials"`
}
