package collector

import (
	"bytes"
	"context"
	"fmt"
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

	cmd := execCommand(ctx, c.command)

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
