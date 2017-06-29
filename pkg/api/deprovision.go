package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *server) deprovision(w http.ResponseWriter, r *http.Request) {
	instanceID := mux.Vars(r)["instance_id"]
	log.Println(instanceID)
	acceptsIncompleteStr := r.URL.Query().Get("accepts_incomplete")
	if acceptsIncompleteStr == "" {
		log.Println("request is missing required query parameter accepts_incomplete=true")
		w.Write(responseAsyncRequired)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	acceptsIncomplete, err := strconv.ParseBool(acceptsIncompleteStr)
	if err != nil || !acceptsIncomplete {
		log.Printf("query paramater accepts_incomplete has invalid value '%s'", acceptsIncompleteStr)
		w.Write(responseAsyncRequired)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	// TODO: Kick off asynchronous deprovisioning
	// TODO: There are actually a lot of different response codes required for
	// different circumstances; also see spec for response body format
	w.WriteHeader(http.StatusAccepted)
}
