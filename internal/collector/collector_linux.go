//go:build linux

package collector

import (
	"time"

	"github.com/dimryb/system-monitor/internal/config"
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

func NewSystemCollector(timeout time.Duration, cfg *config.MonitorConfig) *LinuxCollector {
	metrics := &entity.SystemMetrics{}

	var metricCollectors [metricNumber]metricCollector

	if cfg.Metrics.CPU.Enabled {
		if cfg.Metrics.CPU.CPUUsagePercent {
			metricCollectors[CPUUsagePercent] = &floatMetric{
				value:     &metrics.CPUUsagePercent,
				collector: NewCommandCollector(cpuUsageCommand, timeout),
				parser:    parseFloatMetric,
			}
		}

		if cfg.Metrics.CPU.CPUUserModePercent {
			metricCollectors[CPUUserModePercent] = &floatMetric{
				value:     &metrics.CPUUserModePercent,
				collector: NewCommandCollector(cpuUserModeCommand, timeout),
				parser:    parseFloatMetric,
			}
		}

		if cfg.Metrics.CPU.CPUSystemModePercent {
			metricCollectors[CPUSystemModePercent] = &floatMetric{
				value:     &metrics.CPUSystemModePercent,
				collector: NewCommandCollector(cpuSystemModeCommand, timeout),
				parser:    parseFloatMetric,
			}
		}

		if cfg.Metrics.CPU.CPUIdlePercent {
			metricCollectors[CPUIdlePercent] = &floatMetric{
				value:     &metrics.CPUIdlePercent,
				collector: NewCommandCollector(cpuIdleCommand, timeout),
				parser:    parseFloatMetric,
			}
		}
	}

	if cfg.Metrics.Disk.Enabled {
		if cfg.Metrics.Disk.DiskTPS {
			metricCollectors[DiskTPS] = &floatMetric{
				value:     &metrics.DiskTPS,
				collector: NewCommandCollector(diskIOCommand, timeout),
				parser:    parseDiskTransfersPerSecWithIostat,
			}
		}

		if cfg.Metrics.Disk.DiskKBPerSec {
			metricCollectors[DiskKBPerSec] = &floatMetric{
				value:     &metrics.DiskKBPerSec,
				collector: NewCommandCollector(diskIOCommand, timeout),
				parser:    parseDiskBytesPerSecWithIostat,
			}
		}

		if cfg.Metrics.Disk.DiskUsage {
			metricCollectors[DiskUsage] = &diskUsageMetric{
				value:          &metrics.DiskUsage,
				collectorUsage: NewCommandCollector(diskUsageCommand, timeout),
				collectorInode: NewCommandCollector(diskInodeCommand, timeout),
				parserUsage:    parseDiskUsage,
				parserInode:    parseDiskInodeUsage,
			}
		}
	}

	return &LinuxCollector{
		BaseCollector: BaseCollector{
			timeout:          timeout,
			metrics:          metrics,
			metricCollectors: metricCollectors,
		},
	}
}
