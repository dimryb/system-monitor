package entity

import "time"

type SystemMetrics struct {
	Timestamp time.Time

	CPUUsagePercent      float64
	CPUUserModePercent   float64
	CPUSystemModePercent float64
	CPUIdlePercent       float64

	DiskTPS      float64
	DiskKBPerSec float64

	DiskUsage []DiskUsage
}

type DiskUsage struct {
	Name              string  `json:"name"`
	TotalMB           float64 `json:"total_mb"`
	UsedMB            float64 `json:"used_mb"`
	UsedPercent       float64 `json:"used_percent"`
	InodesTotal       float64 `json:"inodes_total"` // приблизительно
	InodesUsed        float64 `json:"inodes_used"`
	InodesUsedPercent float64 `json:"inodes_used_percent"`
}
