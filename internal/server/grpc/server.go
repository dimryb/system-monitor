package grpc

import (
	"context"
	"fmt"
	"net"

	i "github.com/dimryb/system-monitor/internal/interface"
	"github.com/dimryb/system-monitor/internal/server/grpc/interceptors"
	"github.com/dimryb/system-monitor/proto/monitor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	app i.Application
	cfg ServerConfig
	log i.Logger
}

type ServerConfig struct {
	Port string
}

func NewServer(app i.Application, cfg ServerConfig, log i.Logger) *Server {
	return &Server{
		app: app,
		cfg: cfg,
		log: log,
	}
}

func (s *Server) Run(_ context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.cfg.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.UnaryLoggerInterceptor(s.log)),
	)
	monitor.RegisterSystemMonitorServer(grpcServer, NewMonitorService(s.app))

	reflection.Register(grpcServer)

	s.log.Info("Starting gRPC server", "port", s.cfg.Port)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}
