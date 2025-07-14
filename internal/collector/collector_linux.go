//go:build linux

package collector

import (
	"strconv"
	"strings"
	"time"

	"github.com/dimryb/system-monitor/internal/entity"
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
					parser:    parseFloatMetric,
				},
				CPUUserModePercent: &floatMetric{
					value:     &metrics.CPUUserModePercent,
					collector: NewCommandCollector(cpuUserModeCommand, timeout),
					parser:    parseFloatMetric,
				},
				CPUSystemModePercent: &floatMetric{
					value:     &metrics.CPUSystemModePercent,
					collector: NewCommandCollector(cpuSystemModeCommand, timeout),
					parser:    parseFloatMetric,
				},
				CPUIdlePercent: &floatMetric{
					value:     &metrics.CPUIdlePercent,
					collector: NewCommandCollector(cpuIdleCommand, timeout),
					parser:    parseFloatMetric,
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

func parseFloatMetric(rawData string) (float64, error) {
	str := strings.TrimSpace(rawData)
	load, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return -1.0, err
	}
	return load, nil
}
