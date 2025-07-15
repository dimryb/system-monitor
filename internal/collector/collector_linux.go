//go:build linux

package collector

import (
	"fmt"
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

	diskIOCommand = "iostat -d -k 1 2"

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

func parseNumber(value string) (float64, error) {
	value = strings.ReplaceAll(value, " ", "")

	if strings.Count(value, ".") == 1 {
		parts := strings.Split(value, ".")
		if len(parts) == 2 && len(parts[0]) > 0 && len(parts[1]) > 0 &&
			isNumeric(parts[0]) && isNumeric(parts[1]) {
			return strconv.ParseFloat(value, 64)
		}
	}

	value = strings.ReplaceAll(value, ".", "")

	value = strings.Replace(value, ",", ".", 1)

	return strconv.ParseFloat(value, 64)
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func parseDiskTransfersPerSecWithIostat(rawData string) (float64, error) {
	rawData = strings.ReplaceAll(rawData, "\t", " ")
	lines := strings.Split(rawData, "\n")
	var totalTps float64
	var foundHeader bool

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !foundHeader {
			if strings.HasPrefix(line, "Device ") {
				foundHeader = true
				continue
			}
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		tpsStr := fields[1]
		tps, err := parseNumber(tpsStr)
		if err != nil {
			continue
		}

		totalTps += tps
	}

	if !foundHeader {
		return 0, fmt.Errorf("header not found")
	}

	return totalTps, nil
}

func parseDiskBytesPerSecWithIostat(rawData string) (float64, error) {
	lines := strings.Split(rawData, "\n")
	var totalKB float64
	var foundHeader bool

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if !foundHeader {
			if strings.HasPrefix(line, "Device ") {
				foundHeader = true
			}
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		readKBStr := fields[2]
		writeKBStr := fields[3]

		readKB, _ := parseNumber(readKBStr)
		writeKB, _ := parseNumber(writeKBStr)

		totalKB += readKB + writeKB
	}

	if totalKB == 0 && foundHeader {
		return 0, nil
	}

	if !foundHeader {
		return 0, fmt.Errorf("header not found")
	}

	return totalKB, nil
}
