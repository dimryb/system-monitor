package collector

import (
	"runtime"
	"time"

	i "github.com/dimryb/system-monitor/internal/interface"
)

func NewCollector(timeout time.Duration) i.SystemCollector {
	switch runtime.GOOS {
	case "windows":
		return NewWindowsSystemCollector(timeout)
	case "linux":
		return NewLinuxSystemCollector(timeout)
	default:
		panic("unsupported OS")
	}
}
