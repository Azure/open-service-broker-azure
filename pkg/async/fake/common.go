package fake

import "context"

// RunFunction describes a function used to provide pluggable runtime behavior
// to various fake implementations of interfaces from the async package
type RunFunction func(context.Context) error
