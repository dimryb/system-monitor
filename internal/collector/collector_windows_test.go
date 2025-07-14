//go:build windows

package collector

import (
	"context"
	i "github.com/dimryb/system-monitor/internal/interface"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type fakeCollector struct {
	output string
	err    error
}

func (f *fakeCollector) Collect(ctx context.Context) (string, error) {
	return f.output, f.err
}

func (f *fakeCollector) Timeout() time.Duration {
	return time.Second
}

func TestParseCPULoad(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		wantErr  bool
	}{
		{
			name:     "Valid output with header",
			input:    "LoadPercentage\n\n  75",
			expected: 75,
			wantErr:  false,
		},
		{
			name:     "Only value",
			input:    "  90",
			expected: 90,
			wantErr:  false,
		},
		{
			name:     "Empty input",
			input:    "",
			expected: -1.0,
			wantErr:  true,
		},
		{
			name:     "Invalid number",
			input:    "LoadPercentage\nabc",
			expected: -1.0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := &fakeCollector{output: tt.input}
			result, err := parseCPULoad(context.Background(), collector)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseFloatMetric(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		wantErr  bool
	}{
		{
			name:     "Single valid line",
			input:    "  75.5",
			expected: 75.5,
			wantErr:  false,
		},
		{
			name:     "Multiple lines with empty",
			input:    "\n\n  42\n",
			expected: 42,
			wantErr:  false,
		},
		{
			name:     "No numeric value",
			input:    "invalid data",
			expected: -1.0,
			wantErr:  true,
		},
		{
			name:     "Empty input",
			input:    "",
			expected: -1.0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := &fakeCollector{output: tt.input}
			result, err := parseFloatMetric(context.Background(), collector)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseDiskIO(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		fieldName string
		expected  float64
		wantErr   bool
	}{
		{
			name:      "Valid DiskReadsPersec",
			input:     "DiskReadsPersec : 10",
			fieldName: "DiskReadsPersec",
			expected:  10,
			wantErr:   false,
		},
		{
			name:      "Valid DiskWritesPersec",
			input:     "DiskWritesPersec : 5",
			fieldName: "DiskWritesPersec",
			expected:  5,
			wantErr:   false,
		},
		{
			name:      "Field not found",
			input:     "DiskReadsPersec : 10",
			fieldName: "UnknownField",
			expected:  -1.0,
			wantErr:   true,
		},
		{
			name:      "Malformed line",
			input:     "DiskReadsPersec 10",
			fieldName: "DiskReadsPersec",
			expected:  -1.0,
			wantErr:   true,
		},
		{
			name:      "Empty input",
			input:     "",
			fieldName: "DiskReadsPersec",
			expected:  -1.0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := &fakeCollector{output: tt.input}
			result, err := parseDiskIO(context.Background(), collector, tt.fieldName)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseDiskReadsPerSec(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		wantErr  bool
	}{
		{
			name:     "Valid reads",
			input:    "DiskReadsPersec : 10",
			expected: 10,
			wantErr:  false,
		},
		{
			name:     "No reads field",
			input:    "DiskWritesPersec : 5",
			expected: -1.0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := &fakeCollector{output: tt.input}
			result, err := parseDiskReadsPerSec(context.Background(), collector)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseDiskWritesPerSec(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		wantErr  bool
	}{
		{
			name:     "Valid writes",
			input:    "DiskWritesPersec : 20",
			expected: 20,
			wantErr:  false,
		},
		{
			name:     "No writes field",
			input:    "DiskReadsPersec : 10",
			expected: -1.0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := &fakeCollector{output: tt.input}
			result, err := parseDiskWritesPerSec(context.Background(), collector)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseDiskTransfersPerSec(t *testing.T) {
	tests := []struct {
		name        string
		readOutput  string
		writeOutput string
		expectedTPS float64
		expectError bool
	}{
		{
			name:        "Valid reads and writes",
			readOutput:  "DiskReadsPersec : 10",
			writeOutput: "DiskWritesPersec : 20",
			expectedTPS: 30,
			expectError: false,
		},
		{
			name:        "Empty read output",
			readOutput:  "",
			writeOutput: "DiskWritesPersec : 5",
			expectError: true,
		},
		{
			name:        "Invalid read value",
			readOutput:  "DiskReadsPersec : invalid",
			writeOutput: "DiskWritesPersec : 5",
			expectError: true,
		},
		{
			name:        "Empty write output",
			readOutput:  "DiskReadsPersec : 10",
			writeOutput: "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readCollector := &fakeCollector{output: tt.readOutput}
			writeCollector := &fakeCollector{output: tt.writeOutput}

			result, err := parseDiskTransfersPerSecWithParsers(
				context.Background(),
				nil,
				func(ctx context.Context, _ i.ParamCollector) (float64, error) {
					return parseDiskIO(ctx, readCollector, "DiskReadsPersec")
				},
				func(ctx context.Context, _ i.ParamCollector) (float64, error) {
					return parseDiskIO(ctx, writeCollector, "DiskWritesPersec")
				},
			)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, -1.0, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTPS, result)
			}
		})
	}
}
