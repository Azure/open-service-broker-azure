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
	w.Write(responseEmptyJSON)
	// TODO: Returns 200 or 410; also see spec for response body format
	w.WriteHeader(http.StatusOK)
}
