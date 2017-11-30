package always

import (
	"github.com/Azure/open-service-broker-azure/pkg/api/authenticator"
)

// alwaysAuthenticator is a implementation of the authenticator.Authenticator
// interface useful for testing. It unconditionally authenticates all requests.
type alwaysAuthenticator struct{}

// NewAuthenticator returns an implementation of the authenticator.Authenticator
// interface useful for testing. It unconditionally authenticates all requests.
func NewAuthenticator() authenticator.Authenticator {
	return &alwaysAuthenticator{}
}

func (a *alwaysAuthenticator) Authenticate(
	handler authenticator.HandlerFunction,
) authenticator.HandlerFunction {
	return handler
}
