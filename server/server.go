package server

import (
	"log"
	"net/http"
)

// Server represents the HTTP server for the application
type Server struct {
	router *http.ServeMux
}

// NewServer creates a new server instance
func NewServer() *Server {
	return &Server{
		router: http.NewServeMux(),
	}
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, s.router)
}
