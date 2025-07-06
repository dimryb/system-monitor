package main

import (
	"fmt"

	"github.com/dimryb/system-monitor/internal/app"
	"github.com/dimryb/system-monitor/internal/logger"
)

func main() {
	log := logger.New("debug")

	application := app.NewApp(log)

	fmt.Println(application)

	log.Debugf("Starting system-monitor...")
}
