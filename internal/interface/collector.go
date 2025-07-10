package interfaces

import (
	"context"

	"github.com/dimryb/system-monitor/internal/entity"
)

type SystemCollector interface {
	Collect(ctx context.Context) (*entity.SystemMetrics, error)
}
