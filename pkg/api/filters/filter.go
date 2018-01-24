package filters

import "net/http"

// Filter is a wrapper around a http.HandlerFunc.
type Filter interface {
	Filter(http.HandlerFunc) http.HandlerFunc
}

type chain struct {
	filters []Filter
}

// NewFilterChain returns a Filter capable of applying
// the provided Filter(s) to another http.HandlerFunc
func NewFilterChain(filterChain []Filter) Filter {
	return chain{
		filters: filterChain,
	}
}

// Filter applies a chain of FilterFunction to a target http.HandlerFun
// and returns a new http.HandlerFunc that wraps the chain
func (h chain) Filter(
	target http.HandlerFunc,
) http.HandlerFunc {
	handler := target
	for i := len(h.filters) - 1; i >= 0; i-- {
		filter := h.filters[i]
		handler = filter.Filter(handler)
	}
	return handler
}
