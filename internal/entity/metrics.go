package entity

import "time"

type SystemMetrics struct {
	Timestamp time.Time

	CPUUsagePercent float64
	MemoryUsedMB    uint64
	DiskUsedPercent float64
}
