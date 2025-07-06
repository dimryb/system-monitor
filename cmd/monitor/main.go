package main

import (
	"flag"
	"fmt"
	"log"

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

	cfg, err := config.NewMonitorConfig(configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	logg := logger.New(cfg.Log.Level)

	application := app.NewApp(logg)

	fmt.Println(application)

	logg.Debugf("Starting system-monitor...")
}
