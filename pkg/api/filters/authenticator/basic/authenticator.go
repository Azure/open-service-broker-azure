package basic

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/api/filters"
)

// basicAuthenticator is a implementation of the filter.Filter
// interface that authenticates HTTP requests using Basic Auth
type basicAuthenticator struct {
	Username string
	Password string
}

// NewAuthenticator returns an implementation of the filter.Filter
// interface that authenticates HTTP requests using Basic Auth
func NewAuthenticator(username, password string) filters.Filter {
	return &basicAuthenticator{
		Username: username,
		Password: password,
	}
}

func (b *basicAuthenticator) Filter(
	handle http.HandlerFunc,
) http.HandlerFunc {
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
