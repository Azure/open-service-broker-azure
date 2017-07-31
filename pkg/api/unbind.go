package api

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

func (s *server) unbind(w http.ResponseWriter, r *http.Request) {
	instanceID := mux.Vars(r)["instance_id"]
	log.Debug(instanceID)
	bindingID := mux.Vars(r)["biding_id"]
	log.Debug(bindingID)
	s.writeResponse(w, http.StatusOK, responseEmptyJSON)
}
