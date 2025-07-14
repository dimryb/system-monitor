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
	cpuCollectCommand = `wmic cpu get loadpercentage`

	cpuUserModeCommand          = `(Get-WmiObject -Namespace "root\CIMV2" -Query "SELECT * FROM Win32_PerfFormattedData_Counters_ProcessorInformation WHERE Name='_Total'").PercentUserTime`
	cpuSystemModeCollectCommand = `(Get-WmiObject -Namespace "root\CIMV2" -Query "SELECT * FROM Win32_PerfFormattedData_Counters_ProcessorInformation WHERE Name='_Total'").PercentPrivilegedTime`
	cpuIdleCollectCommand       = `(Get-WmiObject -Namespace "root\CIMV2" -Query "SELECT * FROM Win32_PerfFormattedData_Counters_ProcessorInformation WHERE Name='_Total'").PercentIdleTime`

	diskIOCommand          = `(Get-WmiObject -Namespace "root\CIMV2" -Query "SELECT * FROM Win32_PerfFormattedData_PerfDisk_LogicalDisk WHERE Name='_Total'")`
	diskBytesPerSecCommand = `(Get-WmiObject -Namespace "root\CIMV2" -Query "SELECT * FROM Win32_PerfFormattedData_PerfDisk_LogicalDisk WHERE Name='_Total'").DiskBytesPerSec`
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
			},
		},
	}
}

func parseCPULoad(ctx context.Context, collector i.ParamCollector) (float64, error) {
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

func parseFloatMetric(ctx context.Context, collector i.ParamCollector) (float64, error) {
	raw, err := collector.Collect(ctx)
	if err != nil {
		return -1.0, err
	}

	lines := strings.Split(raw, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		val, err := strconv.ParseFloat(line, 64)
		if err != nil {
			return -1.0, err
		}
		return val, nil
	}
	return -1.0, fmt.Errorf("numeric value not found in command output")
}

func parseDiskIO(ctx context.Context, collector i.ParamCollector, fieldName string) (float64, error) {
	raw, err := collector.Collect(ctx)
	if err != nil {
		return -1.0, err
	}

	lines := strings.Split(raw, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.Contains(line, "Name") {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		if key == fieldName {
			parsedVal, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return -1.0, err
			}
			return parsedVal, nil
		}
	}

	return -1.0, fmt.Errorf("field %s not found in disk IO output", fieldName)
}

func parseDiskReadsPerSec(ctx context.Context, collector i.ParamCollector) (float64, error) {
	return parseDiskIO(ctx, collector, "DiskReadsPersec")
}

func parseDiskWritesPerSec(ctx context.Context, collector i.ParamCollector) (float64, error) {
	return parseDiskIO(ctx, collector, "DiskWritesPersec")
}

func parseDiskTransfersPerSec(ctx context.Context, collector i.ParamCollector) (float64, error) {
	return parseDiskTransfersPerSecWithParsers(ctx, collector, parseDiskReadsPerSec, parseDiskWritesPerSec)
}

func parseDiskTransfersPerSecWithParsers(
	ctx context.Context,
	collector i.ParamCollector,
	readParser func(context.Context, i.ParamCollector) (float64, error),
	writeParser func(context.Context, i.ParamCollector) (float64, error),
) (float64, error) {
	read, err := readParser(ctx, collector)
	if err != nil {
		return -1.0, err
	}
	write, err := writeParser(ctx, collector)
	if err != nil {
		return -1.0, err
	}
	return read + write, nil
}
