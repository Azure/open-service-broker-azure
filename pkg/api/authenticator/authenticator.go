package authenticator

import "net/http"

// HandlerFunction is the signature of any function that can handle an HTTP
// request
type HandlerFunction func(http.ResponseWriter, *http.Request)

// Authenticator is an interface to be implemented by any component capable of
// authenticating HTTP requests
type Authenticator interface {
	// Authenticate returns a function that wraps the provided handler with
	// authentication capability
	Authenticate(HandlerFunction) HandlerFunction
}
