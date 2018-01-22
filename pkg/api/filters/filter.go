package filters

import "net/http"

//Filter is a wrapper around a http.HandlerFunc.
type Filter interface {
	Execute(http.HandlerFunc) http.HandlerFunc
}

//Chain represents a chain of Filters that can be applied
//to a http.HandlrFunc.
type Chain interface {
	Filter(http.HandlerFunc) http.HandlerFunc
}

//Chain contains a number of HandlerFunctions
type chain struct {
	filters []Filter
}

//NewFilterChain returns a chain struct populated with
//the Filter(s) provided
func NewFilterChain(filterChain []Filter) Chain {
	return chain{
		filters: filterChain,
	}
}

//Filter applies a chain of FilterFunction to a target http.HandlerFun
//and returns a new http.HandlerFunc that wraps the chain
func (h chain) Filter(
	target http.HandlerFunc,
) http.HandlerFunc {
	handler := target
	for i := len(h.filters) - 1; i >= 0; i-- {
		filter := h.filters[i]
		handler = filter.Execute(handler)
	}
	return handler
}
