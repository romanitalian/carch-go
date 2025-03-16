package http

import (
	"context"
	"net/http"
	"time"

	"github.com/romanitalian/carch-go/internal/pkg/logger"
	"github.com/romanitalian/carch-go/internal/service"
)

type Server struct {
	srv     *http.Server
	handler *Handler
	log     *logger.Logger
}

func NewServer(cfg *Config, services *service.Services, log *logger.Logger) *Server {
	handler := NewHandler(services, log)
	address := cfg.Address + ":" + cfg.Port
	log.Info("Starting HTTP server", map[string]interface{}{"address": address})
	return &Server{
		handler: handler,
		log:     log,
		srv: &http.Server{
			Addr:           address,
			Handler:        handler,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

func (s *Server) Run() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("Shutting down HTTP server", nil)
	return s.srv.Shutdown(ctx)
}
