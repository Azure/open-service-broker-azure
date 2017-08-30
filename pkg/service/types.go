package service

// ProvisioningParameters is an interface to be implemented by module-specific
// types that represent provisioning parameters. This interface doesn't require
// any functions to be implemented. It exists to improve the clarity of function
// signatures and documentation.
type ProvisioningParameters interface{}

// ProvisioningContext is an interface to be implemented by module-specific
// types that represent provisioning context.
type ProvisioningContext interface {
	GetResourceGroupName() string
}

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
