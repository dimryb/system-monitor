package collector

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/dimryb/system-monitor/internal/entity"
	i "github.com/dimryb/system-monitor/internal/interface"
)

const (
	CPUUsagePercent = iota
	CPUUserModePercent
	CPUSystemModePercent
	CPUIdlePercent

	DiskTPS
	DiskKBPerSec

	MemoryUsedMB
	DiskUsedPercent

	metricNumber
)

type metricCollector interface {
	collect(ctx context.Context) error
}

type floatMetric struct {
	value     *float64
	collector i.ParamCollector
	parser    func(context.Context, i.ParamCollector) (float64, error)
}

func (m *floatMetric) collect(ctx context.Context) error {
	val, err := m.parser(ctx, m.collector)
	if err != nil {
		return err
	}
	*m.value = val
	return nil
}

type intMetric struct { //nolint:unused
	value     *int64
	collector i.ParamCollector
	parser    func(context.Context, i.ParamCollector) (int64, error)
}

func (m *intMetric) collect(ctx context.Context) error { //nolint:unused
	val, err := m.parser(ctx, m.collector)
	if err != nil {
		return err
	}
	*m.value = val
	return nil
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
