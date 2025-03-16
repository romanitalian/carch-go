package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/romanitalian/carch-go/internal/pkg/logger"
	"github.com/romanitalian/carch-go/internal/service"
)

type Server struct {
	services *service.Services
	server   *grpc.Server
	addr     string
	log      *logger.Logger
}

func NewServer(addr string, services *service.Services, log *logger.Logger) *Server {
	s := &Server{
		addr:     addr,
		services: services,
		server:   grpc.NewServer(),
		log:      log,
	}

	// Registration of gRPC services
	// pb.RegisterUserServiceServer(s.server, s)
	log.Info("gRPC server initialized", map[string]interface{}{"address": addr})

	return s
}

// Run starts the gRPC server. If a listener is provided, it will use that listener,
// otherwise it will create a new TCP listener on the configured address.
func (s *Server) Run(listener ...net.Listener) error {
	var l net.Listener
	var err error

	if len(listener) > 0 && listener[0] != nil {
		l = listener[0]
	} else {
		l, err = net.Listen("tcp", s.addr)
		if err != nil {
			s.log.Error("Failed to listen", err, map[string]interface{}{"address": s.addr})
			return err
		}
	}

	s.log.Info("gRPC server starting", map[string]interface{}{"address": s.addr})
	return s.server.Serve(l)
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("Shutting down gRPC server", nil)
	done := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	select {
	case <-ctx.Done():
		s.log.Warn("Shutdown timeout exceeded, forcing shutdown", nil)
		s.server.Stop()
		return ctx.Err()
	case <-done:
		s.log.Info("gRPC server stopped gracefully", nil)
		return nil
	}
}
