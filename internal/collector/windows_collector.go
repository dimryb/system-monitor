package collector

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	i "github.com/dimryb/system-monitor/internal/interface"
)

const (
	cpuCollectCommandWindows = `wmic cpu get loadpercentage`
)

type WindowsCollector struct {
	BaseCollector
}

func NewWindowsSystemCollector(timeout time.Duration) *WindowsCollector {
	return &WindowsCollector{
		BaseCollector: BaseCollector{
			timeout: timeout,
			metrics: map[string]metricCollector{
				"CPUUsagePercent": &floatMetric{
					collector: NewCommandCollector(cpuCollectCommandWindows, timeout),
					parser:    parseCPULoadWindows,
				},
			},
		},
	}
}

func parseCPULoadWindows(ctx context.Context, collector i.ParamCollector) (float64, error) {
	raw, err := collector.Collect(ctx)
	if err != nil {
		return -1.0, err
	}

	lines := strings.Split(raw, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "LoadPercentage" {
			continue
		}
		val, err := strconv.ParseFloat(line, 64)
		if err != nil {
			return -1.0, err
		}
		return val, nil
	}
	return -1.0, fmt.Errorf("cpu load not found")
}
