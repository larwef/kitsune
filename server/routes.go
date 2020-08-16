package server

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// GetRouter sets up the routes and eventual middleware.
func (s *Server) GetRouter() http.Handler {
	router := httprouter.New()

	// Message
	router.Handler(http.MethodPost, "/message", s.addMessage())
	router.Handler(http.MethodGet, "/message/:messageId", s.getMessage())

	// Topic
	router.Handler(http.MethodGet, "/topic", s.getTopics())
	router.Handler(http.MethodGet, "/topic/:topicId", s.getTopic())

	return router
}
