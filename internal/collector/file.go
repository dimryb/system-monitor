package collector

import (
	"context"
	"os"
)

type FileCollector struct {
	path string
}

func NewFileCollector(path string) *FileCollector {
	return &FileCollector{path: path}
}

func (f FileCollector) Collect(ctx context.Context) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	data, err := os.ReadFile(f.path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
