package service

// Stability is a type that represents the relative stability of a service
// module
type Stability int

const (
	// StabilityExperimental represents relative stability of the most immature
	// service modules. At this level of stability, we're not even certain we've
	// built the right thing!
	StabilityExperimental Stability = iota
	// StabilityBeta represents relative stability of the moderately immature and
	// semi-experimental service modules
	StabilityBeta
	// StabilityStable represents relative stability of the mature, production-
	// ready service modules
	StabilityStable
)

// StandardProvisioningParameters encapsulates the handful of provisioning
// parameters that are widely required for ANYTHING provisioned in Azure.
type StandardProvisioningParameters struct {
	Location      string            `json:"location"`
	ResourceGroup string            `json:"resourceGroup"`
	Tags          map[string]string `json:"tags"`
}

// ProvisioningParameters is an interface to be implemented by module-specific
// types that represent provisioning parameters. This interface doesn't require
// any functions to be implemented. It exists to improve the clarity of function
// signatures and documentation.
type ProvisioningParameters interface{}

// StandardProvisioningContext encapsulates the small amount of provisioning
// context that is widely required for ANYTHING provisioned in Azure.
type StandardProvisioningContext struct {
	Location      string            `json:"location"`
	ResourceGroup string            `json:"resourceGroup"`
	Tags          map[string]string `json:"tags"`
}

// ProvisioningContext is an interface to be implemented by module-specific
// types that represent provisioning context.
type ProvisioningContext interface{}

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

// BindingContext is an interface to be implemented by module-specific types
// that represent binding context. This interface doesn't require any functions
// to be implemented. It exists to improve the clarity of function signatures
// and documentation.
type BindingContext interface{}

// Credentials is an interface to be implemented by module-specific types
// that represent service credentials. This interface doesn't require any
// functions to be implemented. It exists to improve the clarity of function
// signatures and documentation.
type Credentials interface{}
