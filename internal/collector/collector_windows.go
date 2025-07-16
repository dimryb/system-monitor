//go:build windows

package collector

import (
	"time"

	"github.com/dimryb/system-monitor/internal/entity"
)

const (
	cpuCollectCommand = `wmic cpu get loadpercentage`

	cpuUserModeCommand          = `(Get-WmiObject -Namespace "root\CIMV2" -Query "SELECT * FROM Win32_PerfFormattedData_Counters_ProcessorInformation WHERE Name='_Total'").PercentUserTime`
	cpuSystemModeCollectCommand = `(Get-WmiObject -Namespace "root\CIMV2" -Query "SELECT * FROM Win32_PerfFormattedData_Counters_ProcessorInformation WHERE Name='_Total'").PercentPrivilegedTime`
	cpuIdleCollectCommand       = `(Get-WmiObject -Namespace "root\CIMV2" -Query "SELECT * FROM Win32_PerfFormattedData_Counters_ProcessorInformation WHERE Name='_Total'").PercentIdleTime`

	diskIOCommand          = `(Get-WmiObject -Namespace "root\CIMV2" -Query "SELECT * FROM Win32_PerfFormattedData_PerfDisk_LogicalDisk WHERE Name='_Total'")`
	diskBytesPerSecCommand = `(Get-WmiObject -Namespace "root\CIMV2" -Query "SELECT * FROM Win32_PerfFormattedData_PerfDisk_LogicalDisk WHERE Name='_Total'").DiskBytesPerSec`

	diskUsageCommand = `(Get-WmiObject -Query "SELECT * FROM Win32_LogicalDisk WHERE DriveType=3") | ConvertTo-Json`
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
					collector: NewCommandCollector(cpuCollectCommand, timeout),
					parser:    parseCPULoad,
				},
				CPUUserModePercent: &floatMetric{
					value:     &metrics.CPUUserModePercent,
					collector: NewCommandCollector(cpuUserModeCommand, timeout),
					parser:    parseFloatMetric,
				},
				CPUSystemModePercent: &floatMetric{
					value:     &metrics.CPUSystemModePercent,
					collector: NewCommandCollector(cpuSystemModeCollectCommand, timeout),
					parser:    parseFloatMetric,
				},
				CPUIdlePercent: &floatMetric{
					value:     &metrics.CPUIdlePercent,
					collector: NewCommandCollector(cpuIdleCollectCommand, timeout),
					parser:    parseFloatMetric,
				},

				DiskTPS: &floatMetric{
					value:     &metrics.DiskTPS,
					collector: NewCommandCollector(diskIOCommand, timeout),
					parser:    parseDiskTransfersPerSec,
				},
				DiskKBPerSec: &floatMetric{
					value:     &metrics.DiskKBPerSec,
					collector: NewCommandCollector(diskBytesPerSecCommand, timeout),
					parser:    parseFloatMetric,
				},

				DiskUsage: &diskUsageMetric{
					value:     &metrics.DiskUsage,
					collector: NewCommandCollector(diskUsageCommand, timeout),
					parser:    parseDiskUsage,
				},
			},
		},
	}
}
