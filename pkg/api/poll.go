package api

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

func (s *server) poll(
	w http.ResponseWriter, // nolint: unparam
	r *http.Request,
) {
	instanceID := mux.Vars(r)["instance_id"]
	log.Debug(instanceID)
	// TODO: Returns 200 or 410; also see spec for response body format
}
