package buffer

import (
	"testing"
	"time"

	"github.com/dimryb/system-monitor/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientBuffer_Add(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name           string
		windowSeconds  int
		initialMetrics []*entity.SystemMetrics
		newMetric      *entity.SystemMetrics
		expectedCount  int
	}{
		{
			name:           "Empty buffer adds one metric",
			windowSeconds:  10,
			initialMetrics: []*entity.SystemMetrics{},
			newMetric: &entity.SystemMetrics{
				CPUUsagePercent: 50.0,
				Timestamp:       now,
			},
			expectedCount: 1,
		},
		{
			name:          "Metric outside window is removed",
			windowSeconds: 5,
			initialMetrics: []*entity.SystemMetrics{
				{
					CPUUsagePercent: 40.0,
					Timestamp:       now.Add(-6 * time.Second),
				},
			},
			newMetric: &entity.SystemMetrics{
				CPUUsagePercent: 60.0,
				Timestamp:       now,
			},
			expectedCount: 1,
		},
		{
			name:          "Multiple metrics inside window are kept",
			windowSeconds: 10,
			initialMetrics: []*entity.SystemMetrics{
				{
					CPUUsagePercent: 30.0,
					Timestamp:       now.Add(-5 * time.Second),
				},
				{
					CPUUsagePercent: 40.0,
					Timestamp:       now.Add(-3 * time.Second),
				},
			},
			newMetric: &entity.SystemMetrics{
				CPUUsagePercent: 50.0,
				Timestamp:       now,
			},
			expectedCount: 3,
		},
		{
			name:          "All initial metrics are out of window",
			windowSeconds: 2,
			initialMetrics: []*entity.SystemMetrics{
				{
					CPUUsagePercent: 30.0,
					Timestamp:       now.Add(-3 * time.Second),
				},
				{
					CPUUsagePercent: 40.0,
					Timestamp:       now.Add(-2 * time.Second),
				},
				{
					CPUUsagePercent: 45.0,
					Timestamp:       now.Add(-1 * time.Second),
				},
			},
			newMetric: &entity.SystemMetrics{
				CPUUsagePercent: 50.0,
				Timestamp:       now,
			},
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := NewClientBuffer(tt.windowSeconds)
			buffer.data = append(buffer.data, tt.initialMetrics...)
			buffer.Add(tt.newMetric)

			assert.Len(t, buffer.data, tt.expectedCount, "The number of elements does not match")

			cutoff := tt.newMetric.Timestamp.Add(-time.Duration(tt.windowSeconds) * time.Second)
			for _, m := range buffer.data {
				require.False(t, m.Timestamp.Before(cutoff), "The %v metric is out of the window", m)
			}
		})
	}
}
