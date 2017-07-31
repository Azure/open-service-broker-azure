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
	w.WriteHeader(statusCode)
	if _, err := w.Write(responseBody); err != nil {
		log.WithField("error", err).Error(
			"api server error: error writing response",
		)
	}
}
