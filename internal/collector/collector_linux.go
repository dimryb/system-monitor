//go:build linux

package collector

import (
	"time"

	"github.com/dimryb/system-monitor/internal/entity"
)

const (
	cpuUsageCommand      = `top -bn1 | grep "Cpu(s)" | awk '{print $2 + $4 + $6}' | sed 's/,/./'`
	cpuUserModeCommand   = `top -bn1 | grep "Cpu(s)" | awk '{print $2}' | sed 's/,/./'`
	cpuSystemModeCommand = `top -bn1 | grep "Cpu(s)" | awk '{print $4}' | sed 's/,/./'`
	cpuIdleCommand       = `top -bn1 | grep "Cpu(s)" | awk '{print $8}' | sed 's/,/./'`

	diskIOCommand = "iostat -d -k 1 2"

	diskUsageCommand = "df -m --output=source,size,used,pcent,target"
	diskInodeCommand = "df -i"
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

				DiskTPS: &floatMetric{
					value:     &metrics.DiskTPS,
					collector: NewCommandCollector(diskIOCommand, timeout),
					parser:    parseDiskTransfersPerSecWithIostat,
				},
				DiskKBPerSec: &floatMetric{
					value:     &metrics.DiskKBPerSec,
					collector: NewCommandCollector(diskIOCommand, timeout),
					parser:    parseDiskBytesPerSecWithIostat,
				},

				DiskUsage: &diskUsageMetric{
					value:          &metrics.DiskUsage,
					collectorUsage: NewCommandCollector(diskUsageCommand, timeout),
					collectorInode: NewCommandCollector(diskInodeCommand, timeout),
					parserUsage:    parseDiskUsage,
					parserInode:    parseDiskInodeUsage,
				},
			},
		},
	}
}
