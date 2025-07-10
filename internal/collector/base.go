package collector

import (
	"context"
	"fmt"
	"time"

	"github.com/dimryb/system-monitor/internal/entity"
	i "github.com/dimryb/system-monitor/internal/interface"
)

const (
	CPUUsagePercent = "CPUUsagePercent"
	MemoryUsedMB    = "MemoryUsedMB"
	DiskUsedPercent = "DiskUsedPercent"
)

type metricCollector interface {
	collect(ctx context.Context) error
	setValue(any)
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

func (m *floatMetric) setValue(value any) {
	val, ok := value.(*float64)
	if !ok {
		panic("not float64")
	}
	m.value = val
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

func (m *intMetric) setValue(value any) { //nolint:unused
	val, ok := value.(*int64)
	if !ok {
		panic("not int64")
	}
	m.value = val
}

type BaseCollector struct {
	metrics map[string]metricCollector
	timeout time.Duration
}

func (c *BaseCollector) Collect(ctx context.Context) (*entity.SystemMetrics, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	metrics := &entity.SystemMetrics{
		Timestamp: time.Now(),
	}

	c.metrics[CPUUsagePercent].setValue(&metrics.CPUUsagePercent)
	//c.metrics[MemoryUsedMB].setValue(&metrics.MemoryUsedMB)
	//c.metrics[DiskUsedPercent].setValue(&metrics.DiskUsedPercent)

	for name, mc := range c.metrics {
		if err := mc.collect(ctx); err != nil {
			return nil, fmt.Errorf("failed to collect metric %q: %w", name, err)
		}
	}

	return metrics, nil
}
