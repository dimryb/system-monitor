package buffer

import (
	"context"
	"github.com/dimryb/system-monitor/internal/config"
	"sync"
	"time"

	"github.com/dimryb/system-monitor/internal/collector"
	i "github.com/dimryb/system-monitor/internal/interface"
)

type GlobalCollector struct {
	collector i.SystemCollector
	buffers   []*ClientBuffer
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	log       i.Logger
}

func NewGlobalCollector(ctx context.Context, log i.Logger, cfg *config.MonitorConfig) *GlobalCollector {
	ctx, cancel := context.WithCancel(ctx)
	return &GlobalCollector{
		collector: collector.NewSystemCollector(2*time.Second, cfg),
		buffers:   make([]*ClientBuffer, 0),
		ctx:       ctx,
		cancel:    cancel,
		log:       log,
	}
}

func (gc *GlobalCollector) Register(windowSeconds int) *ClientBuffer {
	gc.mu.Lock()
	defer gc.mu.Unlock()

	buffer := NewClientBuffer(windowSeconds)
	gc.buffers = append(gc.buffers, buffer)
	return buffer
}

func (gc *GlobalCollector) Start() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-gc.ctx.Done():
			return
		case <-ticker.C:
			metric, err := gc.collector.Collect(gc.ctx)
			if err != nil {
				gc.log.Errorf("failed to collect metric: %v", err)
				continue
			}

			gc.mu.RLock()
			for _, b := range gc.buffers {
				b.Add(metric)
			}
			gc.mu.RUnlock()
		}
	}
}

func (gc *GlobalCollector) Stop() {
	gc.cancel()
}
