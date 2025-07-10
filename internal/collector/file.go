package collector

import (
	"context"
	"os"
)

type fileCollector struct {
	path string
}

func NewFileCollector(path string) *fileCollector {
	return &fileCollector{path: path}
}

func (f fileCollector) Collect(ctx context.Context) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	data, err := os.ReadFile(f.path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
