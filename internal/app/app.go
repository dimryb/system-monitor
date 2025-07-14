package app

import (
	i "github.com/dimryb/system-monitor/internal/interface"
)

type App struct {
	Logger i.Logger
}

func NewApp(logger i.Logger) *App {
	return &App{Logger: logger}
}
