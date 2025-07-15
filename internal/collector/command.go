package collector

import (
	"context"
	"errors"
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

	//fmt.Printf("Executing: %s\n", c.command)

	var out strings.Builder
	var stderr strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	stdout := strings.TrimSpace(out.String())
	stderrOutput := strings.TrimSpace(stderr.String())

	if stderrOutput != "" {
		fmt.Printf("STDERR: %s\n", stderrOutput)
	}

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return "", fmt.Errorf("command timed out: %s", c.command)
	}

	if err != nil {
		return "", fmt.Errorf("command failed: %v, output: %s", err, stdout)
	}

	return stdout, nil
}
