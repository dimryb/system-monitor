package collector

import (
	"context"
	"time"

	"github.com/dimryb/system-monitor/internal/entity"
)

type WindowsCollector struct {
	timeout time.Duration
}

func NewWindowsSystemCollector(timeout time.Duration) *WindowsCollector {
	return &WindowsCollector{timeout: timeout}
}

func (c *WindowsCollector) Collect(ctx context.Context) (*entity.SystemMetrics, error) {
	_ = ctx
	return nil, nil
}
