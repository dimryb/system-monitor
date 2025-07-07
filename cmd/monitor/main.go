package main

import (
	"context"
	"flag"
	"github.com/dimryb/system-monitor/internal/service"
	"log"
	"os/signal"
	"syscall"

	"github.com/dimryb/system-monitor/internal/app"
	"github.com/dimryb/system-monitor/internal/config"
	"github.com/dimryb/system-monitor/internal/logger"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "configs/monitor.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg, err := config.NewMonitorConfig(configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	logg := logger.New(cfg.Log.Level)

	application := app.NewApp(logg)
	monitorService := service.NewMonitorService(application, logg, cfg)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg.Debugf("Starting system-monitor...")
	if err = monitorService.Run(ctx); err != nil {
		logg.Errorf("System-monitor service stopped with error: %v", err)
		cancel()
	} else {
		logg.Infof("System-monitor service stopped gracefully")
	}
}
