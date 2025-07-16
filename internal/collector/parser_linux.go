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
	var result []entity.DiskUsage

	lines := strings.Split(rawData, "\n")
	for _, line := range lines {
		line = sanitizeLine(line)
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		if fields[0] == "Filesystem" {
			continue
		}

		var usePercentStr string
		for i := 3; i < len(fields); i++ {
			if strings.HasSuffix(fields[i], "%") {
				usePercentStr = strings.TrimSuffix(fields[i], "%")
				break
			}
		}

		if usePercentStr == "" {
			continue
		}

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

		result = append(result, entity.DiskUsage{
			Name:        fields[0],
			TotalMB:     totalMB,
			UsedMB:      usedMB,
			UsedPercent: usePercent,
		})
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no disk usage data parsed")
	}

	return result, nil
}

func parseDiskInodeUsage(rawData string) ([]entity.DiskUsage, error) {
	var result []entity.DiskUsage

	lines := strings.Split(rawData, "\n")
	for _, line := range lines {
		line = sanitizeLine(line)
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		if fields[0] == "Filesystem" {
			continue
		}

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

		result = append(result, entity.DiskUsage{
			Name:              fields[0],
			InodesTotal:       inodesTotal,
			InodesUsed:        inodesUsed,
			InodesUsedPercent: ipcent,
		})
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no inode data parsed")
	}

	return result, nil
}

func sanitizeLine(line string) string {
	var builder strings.Builder
	for _, r := range line {
		if r >= 32 && r != 127 {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}
