package service

import (
	"context"
	"sync"

	"github.com/dimryb/system-monitor/internal/config"
	i "github.com/dimryb/system-monitor/internal/interface"
	"github.com/dimryb/system-monitor/internal/server/grpc"
)

type Monitor struct {
	app       i.Application
	logg      i.Logger
	cfg       *config.MonitorConfig
	scheduler *CollectorService
}

func NewMonitorService(ctx context.Context, app i.Application, logger i.Logger, cfg *config.MonitorConfig) *Monitor {
	return &Monitor{
		app:       app,
		logg:      logger,
		cfg:       cfg,
		scheduler: NewCollectorService(ctx, logger),
	}
}

func (m *Monitor) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		m.logg.Debugf("gRPC server starting..")
		grpcServer := grpc.NewServer(
			m.app,
			grpc.ServerConfig{
				Port: m.cfg.GRPC.Port,
			},
			m.logg,
		)
		if err := grpcServer.Run(ctx); err != nil {
			m.logg.Fatalf("Failed to start gRPC server: %m", err.Error())
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := m.scheduler.Start(); err != nil {
			m.logg.Fatalf("Failed to start scheduler: %m", err.Error())
			cancel()
		}
	}()

	m.logg.Infof("System monitor is running...")

	wg.Wait()

	return nil
}
