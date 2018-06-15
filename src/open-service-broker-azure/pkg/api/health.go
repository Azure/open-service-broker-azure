package api

import "net/http"

func (s *server) healthCheck(
	w http.ResponseWriter,
	_ *http.Request,
) {
	if err := s.store.TestConnection(); err != nil {
		s.writeResponse(w, http.StatusInternalServerError, responseEmptyJSON)
		return
	}
	s.writeResponse(w, http.StatusOK, responseEmptyJSON)
}
