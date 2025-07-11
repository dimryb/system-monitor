package collector

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type CommandCollector struct {
	command string
	timeout time.Duration
}

func NewCommandCollector(command string, timeout time.Duration) *CommandCollector {
	return &CommandCollector{command, timeout}
}

func (c CommandCollector) Collect(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.CommandContext(ctx, "powershell.exe", "-Command", c.command) //nolint:gosec
	default:
		cmd = exec.CommandContext(ctx, "bash", "-c", c.command) //nolint:gosec
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	stdout := strings.TrimSpace(out.String())

	if stdout != "" && err != nil {
		return stdout, nil
	}

	if err != nil {
		return "", fmt.Errorf("command failed: %v, output: %s", err, stdout)
	}

	return stdout, nil
}
