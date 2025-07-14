//go:build linux

package collector

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/dimryb/system-monitor/internal/entity"
	i "github.com/dimryb/system-monitor/internal/interface"
)

const (
	cpuUsageCommandLinux      = `top -bn1 | grep "Cpu(s)" | awk '{print $2 + $4 + $6}' | sed 's/,/./'`
	cpuUserModeCommandLinux   = `top -bn1 | grep "Cpu(s)" | awk '{print $2}' | sed 's/,/./'`
	cpuSystemModeCommandLinux = `top -bn1 | grep "Cpu(s)" | awk '{print $4}' | sed 's/,/./'`
	cpuIdleCommandLinux       = `top -bn1 | grep "Cpu(s)" | awk '{print $8}' | sed 's/,/./'`

	//diskCollectCommandLinux = "df -h /"
)

type LinuxCollector struct {
	BaseCollector
}

func NewSystemCollector(timeout time.Duration) *LinuxCollector {
	metrics := &entity.SystemMetrics{}
	return &LinuxCollector{
		BaseCollector: BaseCollector{
			timeout: timeout,
			metrics: metrics,
			metricCollectors: [metricNumber]metricCollector{
				CPUUsagePercent: &floatMetric{
					value:     &metrics.CPUUsagePercent,
					collector: NewCommandCollector(cpuUsageCommandLinux, timeout),
					parser:    parseCPULoadLinux,
				},
				CPUUserModePercent: &floatMetric{
					value:     &metrics.CPUUserModePercent,
					collector: NewCommandCollector(cpuUserModeCommandLinux, timeout),
					parser:    parseCPULoadLinux,
				},
				CPUSystemModePercent: &floatMetric{
					value:     &metrics.CPUSystemModePercent,
					collector: NewCommandCollector(cpuSystemModeCommandLinux, timeout),
					parser:    parseCPULoadLinux,
				},
				CPUIdlePercent: &floatMetric{
					value:     &metrics.CPUIdlePercent,
					collector: NewCommandCollector(cpuIdleCommandLinux, timeout),
					parser:    parseCPULoadLinux,
				},

				//"MemoryUsedMB": intMetric{
				//	collector: NewFileCollector(memoryCollectCommandLinux),
				//	parser:    parseMemoryUsageLinux,
				//},
				//"DiskUsedPercent": floatMetric{
				//	collector: NewCommandCollector(diskCollectCommandLinux, timeout),
				//	parser:    parseDiskUsageLinux,
				//},
			},
		},
	}
}

func parseCPULoadLinux(ctx context.Context, collector i.ParamCollector) (float64, error) {
	raw, err := collector.Collect(ctx)
	if err != nil {
		return -1.0, err
	}

	str := strings.TrimSpace(raw)
	load, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return -1.0, err
	}
	return load, nil
}
