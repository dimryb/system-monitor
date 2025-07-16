//go:build windows

package collector

import (
	"context"

	"github.com/dimryb/system-monitor/internal/entity"
	i "github.com/dimryb/system-monitor/internal/interface"
)

type diskUsageMetric struct {
	value          *[]entity.DiskUsage
	collector      i.ParamCollector
	collectorInode i.ParamCollector
	parser         func(rawData string) ([]entity.DiskUsage, error)
}

func (m *diskUsageMetric) collect(ctx context.Context) error {
	rawUsage, err := m.collector.Collect(ctx)
	if err != nil {
		return err
	}

	diskUsage, err := m.parser(rawUsage)
	if err != nil {
		return err
	}

	*m.value = diskUsage
	return nil
}
