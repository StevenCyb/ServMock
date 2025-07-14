package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/StevenCyb/ServMock/pkg/model"
)

const readWriteTimeout = 30 * time.Second
const idleTimeout = 60 * time.Second

type Server struct {
	http.Server
	behaviorSet *model.BehaviorSet
}

// New creates a new Server instance with the specified listen address.
func New(listen string, behaviorSet *model.BehaviorSet) *Server {
	server := &Server{
		Server: http.Server{
			Addr:         listen,
			ReadTimeout:  readWriteTimeout,
			WriteTimeout: readWriteTimeout,
			IdleTimeout:  idleTimeout,
		},
		behaviorSet: behaviorSet,
	}
	server.Handler = http.HandlerFunc(server.handleRequest)

	return server
}

// handleRequest processes incoming HTTP requests and serves mock responses based on the behavior set.
func (s *Server) SetBehaviorSet(behaviorSet *model.BehaviorSet) {
	s.behaviorSet = behaviorSet
}

// handleRequest is the main request handler for the server.
func (s *Server) Start() <-chan error {
	errorChan := make(chan error, 1)

	go func() {
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errorChan <- err
		}
	}()

	return errorChan
}

// Shutdown gracefully stops the server, allowing for ongoing requests to complete.
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.Server.Shutdown(ctx); err != nil {
		return err
	}
	log.Println("Server shutdown gracefully")
	return nil
}
