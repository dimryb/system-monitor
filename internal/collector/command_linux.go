//go:build linux

package collector

import (
	"context"
	"os"
	"os/exec"
)

func execCommand(ctx context.Context, command string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	cmd.Env = append(os.Environ(), "LANG=C")
	return cmd
}
