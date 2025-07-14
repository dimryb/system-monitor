package collector

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"time"
)

type commandCollector struct {
	command string
	timeout time.Duration
}

func NewCommandCollector(command string, timeout time.Duration) *commandCollector {
	return &commandCollector{command, timeout}
}

func (c commandCollector) Collect(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", c.command) //nolint:gosec
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}
