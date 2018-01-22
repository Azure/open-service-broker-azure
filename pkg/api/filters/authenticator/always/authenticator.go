package always

import (
	"net/http"

	"github.com/Azure/open-service-broker-azure/pkg/api/filters"
)

// alwaysAuthenticator is a implementation of the filters.Filter
// interface useful for testing. It unconditionally authenticates all requests.
type alwaysAuthenticator struct{}

// NewAuthenticator returns an implementation of the filters.Filter
// interface useful for testing. It unconditionally authenticates all requests.
func NewAuthenticator() filters.Filter {
	return &alwaysAuthenticator{}
}

func (a *alwaysAuthenticator) Execute(
	handler http.HandlerFunc,
) http.HandlerFunc {
	return handler
}
