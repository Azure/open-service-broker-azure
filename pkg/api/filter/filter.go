package filter

import "net/http"

// Filter is a wrapper around a http.HandlerFunc.
type Filter interface {
	Filter(http.HandlerFunc) http.HandlerFunc
}
