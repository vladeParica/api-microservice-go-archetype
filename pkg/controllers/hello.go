package controllers

import (
	"net/http"
)

func (s *Server) helloHandler(w http.ResponseWriter, r *http.Request) {
	s.JSONResponse(w, http.StatusOK, map[string]string{"message": "Hello World"})
}
