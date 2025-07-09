package collector

import (
	"context"
	"time"

	"github.com/dimryb/system-monitor/internal/entity"
)

type LinuxCollector struct {
	timeout time.Duration
}

func NewLinuxSystemCollector(timeout time.Duration) *LinuxCollector {
	return &LinuxCollector{timeout: timeout}
}

func (c *LinuxCollector) Collect(ctx context.Context) (*entity.SystemMetrics, error) {
	_ = ctx
	return &entity.SystemMetrics{Timestamp: time.Now()}, nil
}
