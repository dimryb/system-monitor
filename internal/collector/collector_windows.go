//go:build windows

package collector

import (
	"time"

	"github.com/dimryb/system-monitor/internal/config"
	"github.com/dimryb/system-monitor/internal/entity"
)

const (
	cpuCollectCommand = `wmic cpu get loadpercentage`

	cpuUserModeCommand = `(
		Get-WmiObject -Namespace "root\CIMV2" ` +
		`-Query "SELECT * FROM Win32_PerfFormattedData_Counters_ProcessorInformation WHERE Name='_Total'"
).PercentUserTime`
	cpuSystemModeCollectCommand = `(
		Get-WmiObject -Namespace "root\CIMV2" ` +
		`-Query "SELECT * FROM Win32_PerfFormattedData_Counters_ProcessorInformation WHERE Name='_Total'"
).PercentPrivilegedTime`
	cpuIdleCollectCommand = `(
		Get-WmiObject -Namespace "root\CIMV2" ` +
		`-Query "SELECT * FROM Win32_PerfFormattedData_Counters_ProcessorInformation WHERE Name='_Total'"
).PercentIdleTime`

	diskIOCommand = `(
		Get-WmiObject -Namespace "root\CIMV2" ` +
		`-Query "SELECT * FROM Win32_PerfFormattedData_PerfDisk_LogicalDisk WHERE Name='_Total'")`
	diskBytesPerSecCommand = `(
		Get-WmiObject -Namespace "root\CIMV2" ` +
		`-Query "SELECT * FROM Win32_PerfFormattedData_PerfDisk_LogicalDisk WHERE Name='_Total'").DiskBytesPerSec`

	diskUsageCommand = `(Get-WmiObject -Query "SELECT * FROM Win32_LogicalDisk WHERE DriveType=3") | ConvertTo-Json`
)

type WindowsCollector struct {
	BaseCollector
}

func NewSystemCollector(timeout time.Duration, cfg *config.MonitorConfig) *WindowsCollector {
	metrics := &entity.SystemMetrics{}

	var metricCollectors [metricNumber]metricCollector

	if cfg.Metrics.CPU.Enabled {
		if cfg.Metrics.CPU.CPUUsagePercent {
			metricCollectors[CPUUsagePercent] = &floatMetric{
				value:     &metrics.CPUUsagePercent,
				collector: NewCommandCollector(cpuCollectCommand, timeout),
				parser:    parseCPULoad,
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
				collector: NewCommandCollector(cpuSystemModeCollectCommand, timeout),
				parser:    parseFloatMetric,
			}
		}

		if cfg.Metrics.CPU.CPUIdlePercent {
			metricCollectors[CPUIdlePercent] = &floatMetric{
				value:     &metrics.CPUIdlePercent,
				collector: NewCommandCollector(cpuIdleCollectCommand, timeout),
				parser:    parseFloatMetric,
			}
		}
	}

	if cfg.Metrics.Disk.Enabled {
		if cfg.Metrics.Disk.DiskTPS {
			metricCollectors[DiskTPS] = &floatMetric{
				value:     &metrics.DiskTPS,
				collector: NewCommandCollector(diskIOCommand, timeout),
				parser:    parseDiskTransfersPerSec,
			}
		}

		if cfg.Metrics.Disk.DiskKBPerSec {
			metricCollectors[DiskKBPerSec] = &floatMetric{
				value:     &metrics.DiskKBPerSec,
				collector: NewCommandCollector(diskBytesPerSecCommand, timeout),
				parser:    parseFloatMetric,
			}
		}

		if cfg.Metrics.Disk.DiskUsage {
			metricCollectors[DiskUsage] = &diskUsageMetric{
				value:     &metrics.DiskUsage,
				collector: NewCommandCollector(diskUsageCommand, timeout),
				parser:    parseDiskUsage,
			}
		}
	}

	return &WindowsCollector{
		BaseCollector: BaseCollector{
			timeout:          timeout,
			metrics:          metrics,
			metricCollectors: metricCollectors,
		},
	}
}
