package filter

import "net/http"

// GetHandlerFn defines functions used to return an HTTP handler that wraps
// another
type GetHandlerFn func(handle http.HandlerFunc) http.HandlerFunc

type genericFilter struct {
	getHandler GetHandlerFn
}

// NewGenericFilter returns a generic implementation of of the filter.Filter
// interface where the functionality is specified by a function passed in as
// an argument
func NewGenericFilter(getHandler GetHandlerFn) Filter {
	return &genericFilter{
		getHandler: getHandler,
	}
}

func (g *genericFilter) GetHandler(handler http.HandlerFunc) http.HandlerFunc {
	return g.getHandler(handler)
}
