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
	cpuUsageCommand      = `top -bn1 | grep "Cpu(s)" | awk '{print $2 + $4 + $6}' | sed 's/,/./'`
	cpuUserModeCommand   = `top -bn1 | grep "Cpu(s)" | awk '{print $2}' | sed 's/,/./'`
	cpuSystemModeCommand = `top -bn1 | grep "Cpu(s)" | awk '{print $4}' | sed 's/,/./'`
	cpuIdleCommand       = `top -bn1 | grep "Cpu(s)" | awk '{print $8}' | sed 's/,/./'`

	//diskCollectCommand = "df -h /"
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
					collector: NewCommandCollector(cpuUsageCommand, timeout),
					parser:    parseCPULoad,
				},
				CPUUserModePercent: &floatMetric{
					value:     &metrics.CPUUserModePercent,
					collector: NewCommandCollector(cpuUserModeCommand, timeout),
					parser:    parseCPULoad,
				},
				CPUSystemModePercent: &floatMetric{
					value:     &metrics.CPUSystemModePercent,
					collector: NewCommandCollector(cpuSystemModeCommand, timeout),
					parser:    parseCPULoad,
				},
				CPUIdlePercent: &floatMetric{
					value:     &metrics.CPUIdlePercent,
					collector: NewCommandCollector(cpuIdleCommand, timeout),
					parser:    parseCPULoad,
				},

				//"MemoryUsedMB": intMetric{
				//	collector: NewFileCollector(memoryCollectCommand),
				//	parser:    parseMemoryUsageLinux,
				//},
				//"DiskUsedPercent": floatMetric{
				//	collector: NewCommandCollector(diskCollectCommand, timeout),
				//	parser:    parseDiskUsageLinux,
				//},
			},
		},
	}
}

func parseCPULoad(ctx context.Context, collector i.ParamCollector) (float64, error) {
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
