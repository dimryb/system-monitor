package service

import (
	"context"
	"sync"

	"github.com/dimryb/system-monitor/internal/config"
	i "github.com/dimryb/system-monitor/internal/interface"
	"github.com/dimryb/system-monitor/internal/server/grpc"
)

type Monitor struct {
	app  i.Application
	logg i.Logger
	cfg  *config.MonitorConfig
}

func NewMonitorService(app i.Application, logger i.Logger, cfg *config.MonitorConfig) *Monitor {
	return &Monitor{
		app:  app,
		logg: logger,
		cfg:  cfg,
	}
}

func (s *Monitor) Run(ctx context.Context) error {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.logg.Debugf("gRPC server starting..")
		grpcServer := grpc.NewServer(
			s.app,
			grpc.ServerConfig{
				Port: s.cfg.GRPC.Port,
			},
			s.logg,
		)
		if err := grpcServer.Run(ctx); err != nil {
			s.logg.Fatalf("Failed to start gRPC server: %s", err.Error())
		}
	}()

	s.logg.Infof("System monitor is running...")

	wg.Wait()

	return nil
}
