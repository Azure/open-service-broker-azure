package filter

import "net/http"

type chain struct {
	filterChain []Filter
}

// Chain represents a chain of Filter that can be applied
// to a http.HandlerFunc
type Chain interface {
	GetHandler(target http.HandlerFunc) http.HandlerFunc
}

// NewChain returns a Chain capable of applying
// the provided Filter(s) to another http.HandlerFunc
func NewChain(filters ...Filter) Chain {
	return chain{
		filterChain: filters,
	}
}

// GetHandler applies a chain of FilterFunction to a target http.HandlerFun
// and returns a new http.HandlerFunc that wraps the chain
func (h chain) GetHandler(
	target http.HandlerFunc,
) http.HandlerFunc {
	handler := target
	for i := len(h.filterChain) - 1; i >= 0; i-- {
		filter := h.filterChain[i]
		handler = filter.Filter(handler)
	}
	return handler
}
