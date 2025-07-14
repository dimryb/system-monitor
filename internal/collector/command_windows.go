//go:build windows

package collector

import (
	"context"
	"os/exec"
)

func execCommand(ctx context.Context, command string) *exec.Cmd {
	return exec.CommandContext(ctx, "powershell.exe", "-Command", command) //nolint:gosec
}
