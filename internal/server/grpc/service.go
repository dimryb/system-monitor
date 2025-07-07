package grpc

import (
	i "github.com/dimryb/system-monitor/internal/interface"
	"github.com/dimryb/system-monitor/proto/monitor"
)

type MonitorService struct {
	monitor.UnimplementedSystemMonitorServer
	app i.Application
}

func NewMonitorService(app i.Application) *MonitorService {
	return &MonitorService{
		app: app,
	}
}
