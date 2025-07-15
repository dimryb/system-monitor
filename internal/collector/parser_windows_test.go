//go:build windows

package collector

import (
	"encoding/json"
	"github.com/dimryb/system-monitor/internal/entity"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCPULoad(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		wantErr  bool
	}{
		{
			name:     "Valid output with header",
			input:    "LoadPercentage\n\n75",
			expected: 75,
			wantErr:  false,
		},
		{
			name:     "Only value",
			input:    "90",
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
			result, err := parseCPULoad(tt.input)

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
			input:    "75.5",
			expected: 75.5,
			wantErr:  false,
		},
		{
			name:     "Multiple lines with empty",
			input:    "\n\n42\n",
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
			result, err := parseFloatMetric(tt.input)

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
			result, err := parseDiskIO(tt.input, tt.fieldName)

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
		{
			name:     "Invalid value",
			input:    "DiskReadsPersec : invalid",
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
			result, err := parseDiskReadsPerSec(tt.input)

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
		{
			name:     "Invalid value",
			input:    "DiskWritesPersec : invalid",
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
			result, err := parseDiskWritesPerSec(tt.input)

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
		input       string
		expectedTPS float64
		expectError bool
	}{
		{
			name:        "Valid reads and writes",
			input:       "DiskReadsPersec : 10\nDiskWritesPersec : 20",
			expectedTPS: 30,
			expectError: false,
		},
		{
			name:        "Missing reads",
			input:       "DiskWritesPersec : 5",
			expectError: true,
		},
		{
			name:        "Missing writes",
			input:       "DiskReadsPersec : 10",
			expectError: true,
		},
		{
			name:        "Invalid read value",
			input:       "DiskReadsPersec : invalid\nDiskWritesPersec : 5",
			expectError: true,
		},
		{
			name:        "Invalid write value",
			input:       "DiskReadsPersec : 10\nDiskWritesPersec : invalid",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseDiskTransfersPerSec(tt.input)

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

func makeMockDiskUsageJSON(disks []map[string]interface{}) string {
	data, _ := json.Marshal(disks)
	return string(data)
}

func TestParserDiskUsage(t *testing.T) {
	const (
		totalC = 500_000_000_000
		freeC  = 100_000_000_000
		usedC  = totalC - freeC

		totalD = 1_000_000_000_000
		freeD  = 500_000_000_000
		usedD  = totalD - freeD
	)
	tests := []struct {
		name     string
		input    string
		expected []entity.DiskUsage
		wantErr  bool
	}{
		{
			name: "Single disk",
			input: makeMockDiskUsageJSON([]map[string]interface{}{
				{
					"Name":      "C:",
					"Size":      strconv.FormatInt(totalC, 10),
					"FreeSpace": strconv.FormatInt(freeC, 10),
				},
			}),
			expected: []entity.DiskUsage{
				{
					Name:              "C:",
					TotalMB:           float64(totalC) / (1024 * 1024),
					UsedMB:            float64(totalC-freeC) / (1024 * 1024),
					UsedPercent:       float64(usedC*100) / float64(totalC),
					InodesTotal:       float64(totalC) / 4096,
					InodesUsed:        (float64(totalC) / 4096) * 0.8,
					InodesUsedPercent: 80.0,
				},
			},
			wantErr: false,
		},
		{
			name: "Multiple disks",
			input: makeMockDiskUsageJSON([]map[string]interface{}{
				{
					"Name":      "C:",
					"Size":      strconv.FormatInt(totalC, 10),
					"FreeSpace": strconv.FormatInt(freeC, 10),
				},
				{
					"Name":      "D:",
					"Size":      strconv.FormatInt(totalD, 10),
					"FreeSpace": strconv.FormatInt(freeD, 10),
				},
			}),
			expected: []entity.DiskUsage{
				{
					Name:              "C:",
					TotalMB:           float64(totalC) / (1024 * 1024),
					UsedMB:            float64(usedC) / (1024 * 1024),
					UsedPercent:       float64(usedC*100) / float64(totalC),
					InodesTotal:       float64(totalC) / 4096,
					InodesUsed:        (float64(totalC) / 4096) * 0.8,
					InodesUsedPercent: 80.0,
				},
				{
					Name:              "D:",
					TotalMB:           float64(totalD) / (1024 * 1024),
					UsedMB:            float64(usedD) / (1024 * 1024),
					UsedPercent:       float64(usedD*100) / float64(totalD),
					InodesTotal:       float64(totalD) / 4096,
					InodesUsed:        (float64(totalD) / 4096) * 0.5,
					InodesUsedPercent: 50.0,
				},
			},
			wantErr: false,
		},
		{
			name:     "Empty input",
			input:    "",
			expected: nil,
			wantErr:  true,
		},
		{
			name: "Invalid JSON",
			input: `
{
  "Name": "C:",
  "Size": "invalid"
}
`,
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parserDiskUsage(tt.input)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Len(t, result, len(tt.expected))

			for i := range tt.expected {
				assert.Equal(t, tt.expected[i].Name, result[i].Name)
				assert.InDelta(t, tt.expected[i].TotalMB, result[i].TotalMB, 0.1)
				assert.InDelta(t, tt.expected[i].UsedMB, result[i].UsedMB, 0.1)
				assert.InDelta(t, tt.expected[i].UsedPercent, result[i].UsedPercent, 0.1)
				assert.InDelta(t, tt.expected[i].InodesTotal, result[i].InodesTotal, 1)
				assert.InDelta(t, tt.expected[i].InodesUsed, result[i].InodesUsed, 1)
				assert.InDelta(t, tt.expected[i].InodesUsedPercent, result[i].InodesUsedPercent, 0.1)
			}
		})
	}
}
