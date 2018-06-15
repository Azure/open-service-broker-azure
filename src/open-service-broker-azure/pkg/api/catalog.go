package api

import "net/http"

func (s *server) getCatalog(
	w http.ResponseWriter,
	_ *http.Request,
) {
	s.writeResponse(w, http.StatusOK, s.catalogResponse)
}
