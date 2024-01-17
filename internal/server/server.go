package server

import (
	"context"
	"net/http"
	"time"

	"github.com/HeadGardener/effective_mobile/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(conf config.ServerConfig, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + conf.Port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    conf.ReadTimeout * time.Second,
		WriteTimeout:   conf.WriteTimeout * time.Second,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
