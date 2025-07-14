package interfaces

import (
	"context"

	"github.com/dimryb/system-monitor/internal/entity"
)

type Collector interface {
	Collect(ctx context.Context) (*entity.SystemMetrics, error)
}
