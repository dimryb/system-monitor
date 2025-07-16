//go:build linux

package collector

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dimryb/system-monitor/internal/entity"
)

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

func parseDiskUsage(rawData string) ([]entity.DiskUsage, error) {
	entries := parseLines(rawData, 5, func(fields []string) bool {
		_, err := strconv.ParseFloat(fields[1], 64)
		return err == nil
	})

	var disks []entity.DiskUsage
	for _, fields := range entries {
		usePercentStr := strings.TrimSuffix(fields[4], "%")
		usePercent, err := strconv.ParseFloat(usePercentStr, 64)
		if err != nil {
			continue
		}

		totalMB, err := strconv.ParseFloat(fields[1], 64)
		if err != nil {
			continue
		}

		usedMB, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			continue
		}

		disks = append(disks, entity.DiskUsage{
			Name:        fields[0],
			TotalMB:     totalMB,
			UsedMB:      usedMB,
			UsedPercent: usePercent,
		})
	}

	if len(disks) == 0 {
		return nil, fmt.Errorf("no disk usage data parsed")
	}

	return disks, nil
}

func parseDiskInodeUsage(rawData string) ([]entity.DiskUsage, error) {
	entries := parseLines(rawData, 5, func(fields []string) bool {
		_, err := strconv.ParseUint(fields[1], 10, 64)
		return err == nil
	})

	var disks []entity.DiskUsage
	for _, fields := range entries {
		ipcentStr := strings.TrimSuffix(fields[4], "%")
		ipcent, err := strconv.ParseFloat(ipcentStr, 64)
		if err != nil {
			continue
		}

		inodesTotal, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}

		inodesUsed, err := strconv.ParseUint(fields[2], 10, 64)
		if err != nil {
			continue
		}

		disks = append(disks, entity.DiskUsage{
			Name:              fields[0],
			InodesTotal:       inodesTotal,
			InodesUsed:        inodesUsed,
			InodesUsedPercent: ipcent,
		})
	}

	if len(disks) == 0 {
		return nil, fmt.Errorf("no inode data parsed")
	}

	return disks, nil
}

func parseLines(rawData string, minFields int, validLine func(fields []string) bool) [][]string {
	var result [][]string
	lines := strings.Split(rawData, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < minFields {
			continue
		}
		if validLine != nil && !validLine(fields) {
			continue
		}
		result = append(result, fields)
	}
	return result
}
