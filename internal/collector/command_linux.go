//go:build linux

package collector

import (
	"context"
	"os/exec"
)

func execCommand(ctx context.Context, command string) *exec.Cmd {
	return exec.CommandContext(ctx, "bash", "-c", command) //nolint:gosec
}
