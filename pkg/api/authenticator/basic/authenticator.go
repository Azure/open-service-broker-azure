package basic

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/Azure/azure-service-broker/pkg/api/authenticator"
)

// basicAuthenticator is a implementation of the authenticator.Authenticator
// interface useful for testing. It unconditionally authenticates all requests.
type basicAuthenticator struct {
	Username string
	Password string
}

// NewAuthenticator returns an implementation of the authenticator.Authenticator
// interface that authenticates HTTP requests using Basic Auth
func NewAuthenticator(username, password string) authenticator.Authenticator {
	return &basicAuthenticator{
		Username: username,
		Password: password,
	}
}

func (b *basicAuthenticator) Authenticate(
	handle authenticator.HandlerFunction,
) authenticator.HandlerFunction {
	return func(w http.ResponseWriter, r *http.Request) {
		headerValue := r.Header.Get("Authorization")
		if headerValue == "" {
			http.Error(w, "{}", http.StatusUnauthorized)
			return
		}
		headerValueTokens := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(headerValueTokens) != 2 || headerValueTokens[0] != "Basic" {
			http.Error(w, "{}", http.StatusUnauthorized)
			return
		}
		b64UsernameAndPassword := headerValueTokens[1]
		usernameAndPassword, err := base64.StdEncoding.DecodeString(
			b64UsernameAndPassword,
		)
		if err != nil {
			http.Error(w, "{}", http.StatusUnauthorized)
			return
		}
		usernameAndPasswordTokens := strings.SplitN(
			string(usernameAndPassword),
			":",
			2,
		)
		if len(usernameAndPasswordTokens) != 2 ||
			usernameAndPasswordTokens[0] != b.Username ||
			usernameAndPasswordTokens[1] != b.Password {
			http.Error(w, "{}", http.StatusUnauthorized)
			return
		}
		handle(w, r)
	}
}
