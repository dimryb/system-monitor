//go:build windows

package collector

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/dimryb/system-monitor/internal/entity"
)

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

func parseDiskUsage(rawData string) ([]entity.DiskUsage, error) {
	rawData = strings.TrimSpace(rawData)
	if rawData == "" {
		return nil, fmt.Errorf("empty disk usage data")
	}

	var disks []map[string]interface{}
	if err := json.Unmarshal([]byte(rawData), &disks); err != nil {
		return nil, fmt.Errorf("failed to parse disk usage JSON: %w", err)
	}

	var result []entity.DiskUsage

	for _, d := range disks {
		name, ok := d["Name"].(string)
		if !ok || name == "" {
			continue
		}

		total, _ := parseFloatFromInterface(d["Size"])
		free, _ := parseFloatFromInterface(d["FreeSpace"])

		used := (total - free) / (1024 * 1024) // МБ
		totalMB := total / (1024 * 1024)

		usedPercent := 0.0
		if total > 0 {
			usedPercent = used / totalMB * 100
		}

		result = append(result, entity.DiskUsage{
			Name:        name,
			TotalMB:     totalMB,
			UsedMB:      used,
			UsedPercent: usedPercent,
		})
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no valid disk info found")
	}

	return result, nil
}

func parseFloatFromInterface(v interface{}) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case int64:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	default:
		return 0, fmt.Errorf("unsupported type: %T", v)
	}
}
