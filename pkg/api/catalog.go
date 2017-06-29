package api

import "net/http"

func (s *server) getCatalog(w http.ResponseWriter, r *http.Request) {
	w.Write(s.catalogResponse)
}
