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

const (
	// MigrationTag is the tag of migration services. It can be used for tag
	// filter to filter out migration services.
	MigrationTag string = "Migration"
)

// ProvisioningParameters wraps a map containing provisioning parameters.
type ProvisioningParameters struct {
	Parameters
}

// InstanceDetails is an alias for the emoty interface. It exists only to
// improve the clarity of function signatures and documentation.
type InstanceDetails interface{}

// BindingParameters wraps a map containing binding parameters.
type BindingParameters struct {
	Parameters
}

// BindingDetails is an alias for the empty interface. It exists only to improve
// the clarity of function signatures and documentation.
type BindingDetails interface{}

// Credentials is an interface to be implemented by service-specific types
// that represent service credentials. This interface doesn't require any
// functions to be implemented. It exists to improve the clarity of function
// signatures and documentation.
type Credentials interface{}
