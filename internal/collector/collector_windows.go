//go:build windows

package collector

import (
	"fmt"
	"strconv"
	"strings"
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

func parseCPULoad(rawData string) (float64, error) {
	lines := strings.Split(rawData, "\n")
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

func parseFloatMetric(rawData string) (float64, error) {
	lines := strings.Split(rawData, "\n")
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

func parseDiskIO(rawData string, fieldName string) (float64, error) {
	lines := strings.Split(rawData, "\n")
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

func parseDiskReadsPerSec(rawData string) (float64, error) {
	return parseDiskIO(rawData, "DiskReadsPersec")
}

func parseDiskWritesPerSec(rawData string) (float64, error) {
	return parseDiskIO(rawData, "DiskWritesPersec")
}

func parseDiskTransfersPerSec(rawData string) (float64, error) {
	return parseDiskTransfersPerSecWithParsers(rawData, parseDiskReadsPerSec, parseDiskWritesPerSec)
}

func parseDiskTransfersPerSecWithParsers(
	rawData string,
	readParser func(string) (float64, error),
	writeParser func(string) (float64, error),
) (float64, error) {
	read, err := readParser(rawData)
	if err != nil {
		return -1.0, err
	}
	write, err := writeParser(rawData)
	if err != nil {
		return -1.0, err
	}
	return read + write, nil
}
