package headers

import (
	"net/http"
	"strings"

	"github.com/Azure/open-service-broker-azure/pkg/api/filters"
	log "github.com/Sirupsen/logrus"
)

var responseMissingAPIVersion = []byte(
	`{ "error": "MissingAPIVersion", "description": "The request did not ` +
		`include the X-Broker-API-Version header"}`,
)

var responseAPIVersionIncorrect = []byte(
	`{ "error": "APIVersionIncorrect", "description": "X-Broker-API-Verson ` +
		`header includes an incompatible version"}`,
)

type validator struct{}

//NewValidator creates a new instance of the header validator
func NewValidator() filters.Filter {
	return &validator{}
}

//Filter validates that required headers are present
func (v *validator) Filter(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiVersion := r.Header.Get("X-Broker-API-Version")
		if apiVersion == "" {
			v.sendError(w, responseMissingAPIVersion)
			return
		}
		if !strings.HasPrefix(apiVersion, "2.") {
			v.sendError(w, responseAPIVersionIncorrect)
			return
		}
		//call the original handler.
		handler(w, r)
	}
}

func (v *validator) sendError(
	w http.ResponseWriter,
	responseBody []byte,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusPreconditionFailed)
	if _, err := w.Write(responseBody); err != nil {
		log.WithField("error", err).Error(
			"filter error: error writing response",
		)
	}
}
