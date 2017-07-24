package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *server) deprovision(w http.ResponseWriter, r *http.Request) {
	// This broker provisions everything asynchronously. If a client doesn't
	// explicitly indicate that they will accept an incomplete result, the
	// spec says to respond with a 422
	acceptsIncompleteStr := r.URL.Query().Get("accepts_incomplete")
	if acceptsIncompleteStr == "" {
		log.Println(
			"request is missing required query parameter accepts_incomplete=true",
		)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(responseAsyncRequired)
		return
	}
	acceptsIncomplete, err := strconv.ParseBool(acceptsIncompleteStr)
	if err != nil || !acceptsIncomplete {
		log.Printf(
			"query paramater accepts_incomplete has invalid value '%s'",
			acceptsIncompleteStr,
		)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(responseAsyncRequired)
		return
	}
	instanceID := mux.Vars(r)["instance_id"]
	log.Println(instanceID)
	// TODO: Kick off asynchronous deprovisioning
	// TODO: There are actually a lot of different response codes required for
	// different circumstances; also see spec for response body format
	w.WriteHeader(http.StatusAccepted)
}
