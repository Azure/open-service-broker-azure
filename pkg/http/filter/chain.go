package filter

import "net/http"

// NewChain returns a Filter that can wrap a succession of http.HandlerFuncs
// provided by other Filters around another http.HandlerFunc
func NewChain(filters ...Filter) Filter {
	return NewGenericFilter(
		func(handler http.HandlerFunc) http.HandlerFunc {
			for i := len(filters) - 1; i >= 0; i-- {
				handler = filters[i].GetHandler(handler)
			}
			return handler
		},
	)
}
