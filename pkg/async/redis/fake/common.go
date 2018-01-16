package fake

import "context"

// RunFn describes a function used to provide pluggable runtime behavior to
// various fake implementations of interfaces from the async package
type RunFn func(context.Context) error
