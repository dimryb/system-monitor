package collector

import (
	"context"

	i "github.com/dimryb/system-monitor/internal/interface"
)

type floatMetric struct {
	value     *float64
	collector i.ParamCollector
	parser    func(rawData string) (float64, error)
}

func (m *floatMetric) collect(ctx context.Context) error {
	raw, err := m.collector.Collect(ctx)
	if err != nil {
		return err
	}

	val, err := m.parser(raw)
	if err != nil {
		return err
	}
	*m.value = val
	return nil
}
