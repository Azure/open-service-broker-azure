package filters

import (
	"net/http"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/http/filter"
	log "github.com/Sirupsen/logrus"
)

var (
	responseMissingAPIVersion = []byte(
		`{ "error": "MissingAPIVersion", "description": "The request did not ` +
			`include the X-Broker-API-Version header"}`,
	)
	responseAPIVersionIncorrect = []byte(
		`{ "error": "APIVersionIncorrect", "description": "X-Broker-API-Verson ` +
			`header includes an incompatible version"}`,
	)
)

// NewAPIVersionFilter returns an implementation of the filter.Filter interface
// that validates that an API version is specified by a request's
// X-Broker-API-Version header and that the specified version is compatible
// with this broker.
func NewAPIVersionFilter() filter.Filter {
	return filter.NewGenericFilter(
		func(handle http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				apiVersion := r.Header.Get("X-Broker-API-Version")
				if apiVersion == "" {
					sendError(w, responseMissingAPIVersion)
					return
				}
				if !strings.HasPrefix(apiVersion, "2.") {
					sendError(w, responseAPIVersionIncorrect)
					return
				}
				// Call the original handler
				handle(w, r)
			}
		},
	)
}

func sendError(w http.ResponseWriter, responseBody []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusPreconditionFailed)
	if _, err := w.Write(responseBody); err != nil {
		log.WithField("error", err).Error(
			"filter error: error writing response",
		)
	}
}
