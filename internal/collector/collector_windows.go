//go:build windows

package collector

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dimryb/system-monitor/internal/entity"
	i "github.com/dimryb/system-monitor/internal/interface"
)

const (
	cpuCollectCommandWindows           = `wmic cpu get loadpercentage`
	cpuUserModeCommandWindows          = `(Get-WmiObject -Namespace "root\CIMV2" -Query "SELECT * FROM Win32_PerfFormattedData_Counters_ProcessorInformation WHERE Name='_Total'").PercentUserTime`
	cpuSystemModeCollectCommandWindows = `(Get-WmiObject -Namespace "root\CIMV2" -Query "SELECT * FROM Win32_PerfFormattedData_Counters_ProcessorInformation WHERE Name='_Total'").PercentPrivilegedTime`
	cpuIdleCollectCommandWindows       = `(Get-WmiObject -Namespace "root\CIMV2" -Query "SELECT * FROM Win32_PerfFormattedData_Counters_ProcessorInformation WHERE Name='_Total'").PercentIdleTime`
)

type WindowsCollector struct {
	BaseCollector
}

func NewSystemCollector(timeout time.Duration) *WindowsCollector {
	metrics := &entity.SystemMetrics{}
	return &WindowsCollector{
		BaseCollector: BaseCollector{
			timeout: timeout,
			metrics: metrics,
			metricCollectors: [metricNumber]metricCollector{
				CPUUsagePercent: &floatMetric{
					value:     &metrics.CPUUsagePercent,
					collector: NewCommandCollector(cpuCollectCommandWindows, timeout),
					parser:    parseCPULoadWindows,
				},
				CPUUserModePercent: &floatMetric{
					value:     &metrics.CPUUserModePercent,
					collector: NewCommandCollector(cpuUserModeCommandWindows, timeout),
					parser:    parseCPULoadWindows,
				},
				CPUSystemModePercent: &floatMetric{
					value:     &metrics.CPUSystemModePercent,
					collector: NewCommandCollector(cpuSystemModeCollectCommandWindows, timeout),
					parser:    parseCPULoadWindows,
				},
				CPUIdlePercent: &floatMetric{
					value:     &metrics.CPUIdlePercent,
					collector: NewCommandCollector(cpuIdleCollectCommandWindows, timeout),
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
