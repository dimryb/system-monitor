package collector

import (
	"context"
	i "github.com/dimryb/system-monitor/internal/interface"
)

type intMetric struct { //nolint:unused
	value     *int64
	collector i.ParamCollector
	parser    func(rawData string) (int64, error)
}

func (m *intMetric) collect(ctx context.Context) error { //nolint:unused
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
