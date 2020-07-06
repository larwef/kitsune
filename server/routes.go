package server

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// GetRouter sets up the routes and eventual middleware.
func (s *Server) GetRouter() http.Handler {
	router := httprouter.New()

	router.Handler(http.MethodPost, "/topic/:topicId", s.publish())
	router.Handler(http.MethodGet, "/topic/:topicId/:messageId", s.getMessage())

	return router
}
