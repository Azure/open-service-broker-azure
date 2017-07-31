package api

import "net/http"

func (s *server) getCatalog(
	w http.ResponseWriter,
	r *http.Request, // nolint: unparam
) {
	s.writeResponse(w, http.StatusOK, s.catalogResponse)
}
