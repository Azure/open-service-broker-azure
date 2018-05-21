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

// ProvisioningParameters is an alias for maps intended to contain non-sensitive
// provisioning parameters. It exists only to improve the clarity of function
// signatures and documentation.
type ProvisioningParameters map[string]interface{}

// SecureProvisioningParameters is an alias for maps intended to contain
// sensitive provisioning parameters. It exists only to improve the clarity of
// function signatures and documentation.
type SecureProvisioningParameters map[string]interface{}

// InstanceDetails is an alias for maps intended to contain non-sensitive
// details of a service instance. It exists only to improve the clarity of
// function signatures and documentation.
type InstanceDetails map[string]interface{}

// SecureInstanceDetails is an alias for maps intended to contain sensitive
// details of a service instance. It exists only to improve the clarity of
// function signatures and documentation.
type SecureInstanceDetails map[string]interface{}

// CombinedBindingParameters is an alias for maps intended to contain inbound
// binding parameters-- which may contain both sensitive and non-sensitive
// values. It exists only to improve the clarity of function signatures and
// documentation.
type CombinedBindingParameters map[string]interface{}

// BindingParameters is an alias for maps intended to contain non-sensitive
// binding parameters. It exists only to improve the clarity of function
// signatures and documentation.
type BindingParameters map[string]interface{}

// SecureBindingParameters is an alias for maps intended to contain sensitive
// binding parameters. It exists only to improve the clarity of function
// signatures and documentation.
type SecureBindingParameters map[string]interface{}

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
