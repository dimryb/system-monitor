package entity

import "time"

type SystemMetrics struct {
	Timestamp time.Time

	CPUUsagePercent      float64
	CPUUserModePercent   float64
	CPUSystemModePercent float64
	CPUIdlePercent       float64

	MemoryUsedMB    int64
	DiskUsedPercent float64
}
