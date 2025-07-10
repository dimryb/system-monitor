package entity

import "time"

type SystemMetrics struct {
	Timestamp time.Time

	CPUUsagePercent float64
	MemoryUsedMB    int64
	DiskUsedPercent float64
}
