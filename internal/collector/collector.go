package collector

import (
	"time"

	i "github.com/dimryb/system-monitor/internal/interface"
)

func NewCollector(timeout time.Duration) i.SystemCollector {
	return NewSystemCollector(timeout)
}
