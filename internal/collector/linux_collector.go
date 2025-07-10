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
	cpuCollectCommand  = `top -bn1 | grep "Cpu(s)" | awk '{print $2 + $4}' | sed 's/,/./'`
	diskCollectCommand = "df -h /"
)

type LinuxCollector struct {
	cpuCollector    i.ParamCollector
	memoryCollector i.ParamCollector
	diskCollector   i.ParamCollector
	timeout         time.Duration
}

func NewLinuxSystemCollector(timeout time.Duration) *LinuxCollector {
	return &LinuxCollector{
		cpuCollector:    NewCommandCollector(cpuCollectCommand, timeout),
		memoryCollector: NewFileCollector("/proc/meminfo"),
		diskCollector:   NewCommandCollector(diskCollectCommand, timeout),
		timeout:         timeout,
	}
}

func (c *LinuxCollector) Collect(ctx context.Context) (*entity.SystemMetrics, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	cpuLoad, err := parseCPULoad(ctx, c.cpuCollector)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CPU load: %w", err)
	}

	return &entity.SystemMetrics{
		Timestamp:       time.Now(),
		CPUUsagePercent: cpuLoad,
		MemoryUsedMB:    0,
		DiskUsedPercent: -1,
	}, nil
}

func parseCPULoad(ctx context.Context, collector i.ParamCollector) (float64, error) {
	raw, err := collector.Collect(ctx)
	if err != nil {
		return -1.0, err
	}

	str := strings.TrimSpace(raw)
	load, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return -1.0, err
	}
	return load, nil
}
