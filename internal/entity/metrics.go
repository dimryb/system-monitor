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
	TotalMB           float64 `json:"total_mb"`            //nolint: tagliatelle
	UsedMB            float64 `json:"used_mb"`             //nolint: tagliatelle
	UsedPercent       float64 `json:"used_percent"`        //nolint: tagliatelle
	InodesTotal       uint64  `json:"inodes_total"`        //nolint: tagliatelle
	InodesUsed        uint64  `json:"inodes_used"`         //nolint: tagliatelle
	InodesUsedPercent float64 `json:"inodes_used_percent"` //nolint: tagliatelle
}
