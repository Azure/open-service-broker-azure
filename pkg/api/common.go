package api

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func (s *server) writeResponse(
	w http.ResponseWriter,
	statusCode int,
	responseBody []byte,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(responseBody); err != nil {
		log.WithField("error", err).Error(
			"api server error: error writing response",
		)
	}
}
