package kitsune

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (s *Server) GetRouter() http.Handler {
	router := httprouter.New()

	router.Handler(http.MethodGet, "/", s.Hello())

	return router
}
