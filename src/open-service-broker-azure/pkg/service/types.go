package service

// Stability is a type that represents the relative stability of a service
// module
type Stability int

const (
	// StabilityExperimental represents relative stability of the most immature
	// service modules. At this level of stability, we're not even certain we've
	// built the right thing!
	StabilityExperimental Stability = iota
	// StabilityPreview represents relative stability of modules we believe are
	// approaching a stable state.
	StabilityPreview
	// StabilityStable represents relative stability of the mature, production-
	// ready service modules.
	StabilityStable
)

// ProvisioningParameters wraps a map containing provisioning parameters.
type ProvisioningParameters struct {
	Parameters
}

// InstanceDetails is an alias for maps intended to contain non-sensitive
// details of a service instance. It exists only to improve the clarity of
// function signatures and documentation.
type InstanceDetails map[string]interface{}

// SecureInstanceDetails is an alias for maps intended to contain sensitive
// details of a service instance. It exists only to improve the clarity of
// function signatures and documentation.
type SecureInstanceDetails map[string]interface{}

// BindingParameters wraps a map containing binding parameters.
type BindingParameters struct {
	Parameters
}

// BindingDetails is an alias for maps intended to contain non-sensitive
// details of a service binding. It exists only to improve the clarity of
// function signatures and documentation.
type BindingDetails map[string]interface{}

// SecureBindingDetails is an alias for maps intended to contain sensitive
// details of a service binding. It exists only to improve the clarity of
// function signatures and documentation.
type SecureBindingDetails map[string]interface{}

// Credentials is an interface to be implemented by service-specific types
// that represent service credentials. This interface doesn't require any
// functions to be implemented. It exists to improve the clarity of function
// signatures and documentation.
type Credentials interface{}
