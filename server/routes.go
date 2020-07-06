package server

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// GetRouter sets up the routes and eventual middleware.
func (s *Server) GetRouter() http.Handler {
	router := httprouter.New()

	// Topic
	router.Handler(http.MethodGet, "/topic", s.getTopics())
	router.Handler(http.MethodGet, "/topic/:topicId", s.getTopic())
	router.Handler(http.MethodPost, "/topic/:topicId", s.addMessage())
	router.Handler(http.MethodGet, "/topic/:topicId/:messageId", s.getMessage())

	return router
}
