package http

import (
	"context"
	"errors"
	"net/http"
)

type Server struct {
	srv *http.Server
}

func NewServer(handler http.Handler, addr string) *Server {
	return &Server{
		srv: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *Server) Start() error {
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
