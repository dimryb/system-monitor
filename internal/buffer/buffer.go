package buffer

import (
	"sync"
	"time"

	"github.com/dimryb/system-monitor/internal/entity"
)

type ClientBuffer struct {
	windowSeconds int
	data          []*entity.SystemMetrics
	mu            sync.Mutex
}

func NewClientBuffer(windowSeconds int) *ClientBuffer {
	return &ClientBuffer{
		windowSeconds: windowSeconds,
		data:          make([]*entity.SystemMetrics, 0),
	}
}

func (b *ClientBuffer) Add(metric *entity.SystemMetrics) {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := metric.Timestamp
	b.data = append(b.data, metric)

	cutoff := now.Add(-time.Duration(b.windowSeconds) * time.Second)
	for len(b.data) > 0 && b.data[0].Timestamp.Before(cutoff) {
		b.data = b.data[1:]
	}
}
