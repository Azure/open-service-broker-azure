package always

import (
	"net/http"

	"github.com/Azure/open-service-broker-azure/pkg/api/filter"
)

// alwaysAuthenticator is a implementation of the filter.Filter
// interface useful for testing. It unconditionally authenticates all requests.
type alwaysAuthenticator struct{}

// NewAuthenticator returns an implementation of the filter.Filter
// interface useful for testing. It unconditionally authenticates all requests.
func NewAuthenticator() filter.Filter {
	return &alwaysAuthenticator{}
}

func (a *alwaysAuthenticator) Filter(
	handler http.HandlerFunc,
) http.HandlerFunc {
	return handler
}