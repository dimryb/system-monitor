package collector

import (
	"context"
	"strconv"
	"strings"
	"time"

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

func NewLinuxSystemCollector(timeout time.Duration) *LinuxCollector {
	return &LinuxCollector{
		BaseCollector: BaseCollector{
			timeout: timeout,
			metrics: [metricNumber]metricCollector{
				CPUUsagePercent: &floatMetric{
					collector: NewCommandCollector(cpuUsageCommandLinux, timeout),
					parser:    parseCPULoadLinux,
				},
				CPUUserModePercent: &floatMetric{
					collector: NewCommandCollector(cpuUserModeCommandLinux, timeout),
					parser:    parseCPULoadLinux,
				},
				CPUSystemModePercent: &floatMetric{
					collector: NewCommandCollector(cpuSystemModeCommandLinux, timeout),
					parser:    parseCPULoadLinux,
				},
				CPUIdlePercent: &floatMetric{
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
