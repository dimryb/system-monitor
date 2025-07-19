//go:build linux

package collector

import (
	"context"
	"testing"

	"github.com/dimryb/system-monitor/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockParamCollector struct {
	response string
	err      error
}

func (m *mockParamCollector) Collect(_ context.Context) (string, error) {
	return m.response, m.err
}

func TestDiskUsageMetric_Collect(t *testing.T) {
	rawUsage := `
Filesystem           1M-blocks    Used Available Use% Mounted on
/dev/sdc3              100000     50000     50000   50% /
tmpfs                   2000         1      1999    1% /run
`

	rawInode := `
Filesystem            Inodes   IUsed   IFree IUse% Mounted on
/dev/sdc3            54034432 1400000  100000   3% /
tmpfs                 2040446    2000   10000   1% /run
`

	expected := []entity.DiskUsage{
		{
			Name:              "/dev/sdc3",
			TotalMB:           100000,
			UsedMB:            50000,
			UsedPercent:       50.0,
			InodesTotal:       54034432,
			InodesUsed:        1400000,
			InodesUsedPercent: 3.0,
		},
		{
			Name:              "tmpfs",
			TotalMB:           2000,
			UsedMB:            1,
			UsedPercent:       1.0,
			InodesTotal:       2040446,
			InodesUsed:        2000,
			InodesUsedPercent: 1.0,
		},
	}

	metric := &diskUsageMetric{
		value: new([]entity.DiskUsage),
		collectorUsage: &mockParamCollector{
			response: rawUsage,
		},
		parserUsage: parseDiskUsage,
		collectorInode: &mockParamCollector{
			response: rawInode,
		},
		parserInode: parseDiskInodeUsage,
	}

	err := metric.collect(context.Background())
	require.NoError(t, err)
	require.NotNil(t, *metric.value)
	require.Len(t, *metric.value, len(expected))

	for i := range expected {
		disk := (*metric.value)[i]
		exp := expected[i]

		assert.Equal(t, exp.Name, disk.Name)
		assert.InDelta(t, exp.TotalMB, disk.TotalMB, 0.1)
		assert.InDelta(t, exp.UsedMB, disk.UsedMB, 0.1)
		assert.InDelta(t, exp.UsedPercent, disk.UsedPercent, 0.1)
		assert.InDelta(t, exp.InodesTotal, disk.InodesTotal, 0.1)
		assert.InDelta(t, exp.InodesUsed, disk.InodesUsed, 0.1)
		assert.InDelta(t, exp.InodesUsedPercent, disk.InodesUsedPercent, 0.1)
	}
}
