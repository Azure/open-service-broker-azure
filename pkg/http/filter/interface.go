package filter

import "net/http"

// Filter is an interface to be implemented by components that can wrapper a
// new http.HandlerFunc around another.
type Filter interface {
	// GetHandler decorates one http.HandlerFunc with another
	GetHandler(http.HandlerFunc) http.HandlerFunc
}
