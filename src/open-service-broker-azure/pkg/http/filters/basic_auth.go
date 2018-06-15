package filters

import (
	"encoding/base64"
	"net/http"
	"strings"

	"open-service-broker-azure/pkg/http/filter"
)

// NewBasicAuthFilter returns an implementation of the filter.Filter interface
// that authenticates HTTP requests using Basic Auth
func NewBasicAuthFilter(username, password string) filter.Filter {
	return filter.NewGenericFilter(
		func(handle http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				headerValue := r.Header.Get("Authorization")
				if headerValue == "" {
					http.Error(w, "{}", http.StatusUnauthorized)
					return
				}
				headerValueTokens := strings.SplitN(
					r.Header.Get("Authorization"),
					" ",
					2,
				)
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
					usernameAndPasswordTokens[0] != username ||
					usernameAndPasswordTokens[1] != password {
					http.Error(w, "{}", http.StatusUnauthorized)
					return
				}
				handle(w, r)
			}
		},
	)
}
