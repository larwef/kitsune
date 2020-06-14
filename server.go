package kitsune

import (
	"go.uber.org/zap"
	"net/http"
)

type Server struct{}

func (s *Server) Hello() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		zap.S().Info("Received request")
		res.Header().Set("Content-Type", "text/plain")
		res.Write([]byte("hello"))
	}
}
