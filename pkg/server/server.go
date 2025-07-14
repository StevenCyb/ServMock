package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/StevenCyb/ServMock/pkg/model"
)

type Server struct {
	http.Server
	behaviorSet *model.BehaviorSet
}

func New(listen string) *Server {
	server := &Server{
		Server: http.Server{
			Addr:         listen,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		behaviorSet: &model.BehaviorSet{},
	}
	server.Handler = http.HandlerFunc(server.handleRequest)

	return server
}

func (s *Server) Start() <-chan error {
	errorChan := make(chan error, 1)

	go func() {
		if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errorChan <- err
		}
	}()

	return errorChan
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.Server.Shutdown(ctx); err != nil {
		return err
	}
	log.Println("Server shutdown gracefully")
	return nil
}
