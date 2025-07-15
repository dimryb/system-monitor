package collector

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/dimryb/system-monitor/internal/entity"
)

const (
	CPUUsagePercent = iota
	CPUUserModePercent
	CPUSystemModePercent
	CPUIdlePercent

	DiskTPS
	DiskKBPerSec

	DiskUsage

	MemoryUsedMB
	DiskUsedPercent

	metricNumber
)

type metricCollector interface {
	collect(ctx context.Context) error
}

type BaseCollector struct {
	metricCollectors [metricNumber]metricCollector
	metrics          *entity.SystemMetrics
	timeout          time.Duration
}

func (c *BaseCollector) Collect(ctx context.Context) (*entity.SystemMetrics, error) {
	timestamp := time.Now()
	metrics := entity.SystemMetrics{}
	wg := sync.WaitGroup{}
	errChan := make(chan error, len(c.metricCollectors))

	for ind, mc := range c.metricCollectors {
		if mc == nil {
			continue
		}
		wg.Add(1)
		go func(name int, collector metricCollector) {
			defer wg.Done()
			if err := mc.collect(ctx); err != nil {
				errChan <- fmt.Errorf("failed to collect metric %q: %w", name, err)
			}
		}(ind, mc)
	}

	wg.Wait()
	close(errChan)

	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	metrics = *c.metrics
	metrics.Timestamp = timestamp
	fmt.Println(metrics)

	return &metrics, nil
}
