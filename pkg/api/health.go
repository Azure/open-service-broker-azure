package api

import "net/http"

func (s *server) healthCheck(w http.ResponseWriter, r *http.Request) {
	if err := s.store.TestConnection(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
