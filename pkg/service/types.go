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

// ProvisioningParameters is an interface to be implemented by module-specific
// types that represent provisioning parameters. This interface doesn't require
// any functions to be implemented. It exists to improve the clarity of function
// signatures and documentation.
type ProvisioningParameters interface{}

// InstanceDetails is an interface to be implemented by service-specific
// types that represent the non-sensitive details of a service instance.
type InstanceDetails interface{}

// SecureInstanceDetails is an interface to be implemented by service-specific
// types that represent the secure (sensitive) details of a service instance.
type SecureInstanceDetails interface{}

// UpdatingParameters is an interface to be implemented by module-specific
// types that represent updating parameters. This interface doesn't require
// any functions to be implemented. It exists to improve the clarity of function
// signatures and documentation.
type UpdatingParameters interface{}

// BindingParameters is an interface to be implemented by module-specific types
// that represent binding parameters. This interface doesn't require any
// functions to be implemented. It exists to improve the clarity of function
// signatures and documentation.
type BindingParameters interface{}

// BindingDetails is an interface to be implemented by service-specific types
// that represent non-sensitive binding details. This interface doesn't require
// any functions to be implemented. It exists to improve the clarity of function
// signatures and documentation.
type BindingDetails interface{}

// SecureBindingDetails is an interface to be implemented by service-specific
// types that represent secure (sensitive) binding details. This interface
// doesn't require any functions to be implemented. It exists to improve the
// clarity of function signatures and documentation.
type SecureBindingDetails interface{}

// Credentials is an interface to be implemented by module-specific types
// that represent service credentials. This interface doesn't require any
// functions to be implemented. It exists to improve the clarity of function
// signatures and documentation.
type Credentials interface{}
