package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/dimryb/system-monitor/internal/buffer"
	"github.com/dimryb/system-monitor/internal/config"
	i "github.com/dimryb/system-monitor/internal/interface"
	"github.com/dimryb/system-monitor/internal/server/grpc"
)

type Monitor struct {
	app             i.Application
	log             i.Logger
	cfg             *config.MonitorConfig
	globalCollector *buffer.GlobalCollector
}

func NewMonitorService(ctx context.Context, app i.Application, logger i.Logger, cfg *config.MonitorConfig) *Monitor {
	return &Monitor{
		app:             app,
		log:             logger,
		cfg:             cfg,
		globalCollector: buffer.NewGlobalCollector(ctx, logger),
	}
}

func (m *Monitor) Run(ctx context.Context) error {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		m.log.Debugf("gRPC server starting..")
		grpcServer := grpc.NewServer(
			m.app,
			grpc.ServerConfig{
				Port: m.cfg.GRPC.Port,
			},
			m.log,
		)
		if err := grpcServer.Run(ctx); err != nil {
			m.log.Fatalf("Failed to start gRPC server: %m", err.Error())
		}
	}()

	globalCollector := buffer.NewGlobalCollector(ctx, m.log)
	buf := globalCollector.Register(15)

	wg.Add(1)
	go func() {
		defer wg.Done()
		globalCollector.Start()
	}()

	// Stub register collector
	wg.Add(1)
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				m.log.Infof("Stub collector stopped.")
				return
			case <-ticker.C:
				data := buf.Get()
				m.log.Debugf("Buffer size: %d", len(data))
				fmt.Println("Data:")
				for _, v := range data {
					fmt.Println(*v)
				}
			}
		}
	}()
	// end Stub

	m.log.Infof("System monitor is running...")

	wg.Wait()

	return nil
}
