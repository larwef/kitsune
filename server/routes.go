package server

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (s *Server) GetRouter() http.Handler {
	router := httprouter.New()

	router.Handler(http.MethodPost, "/publish/:topic", s.publish())
	router.Handler(http.MethodPost, "/poll/:topic", s.poll())
	router.Handler(http.MethodGet, "/:topic/:messageId", s.getMessage())

	return router
}