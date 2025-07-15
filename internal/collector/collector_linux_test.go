package collector

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseDiskTransfersPerSecWithIostat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		err      bool
	}{
		{
			name: "One disk",
			input: `Linux 5.15.0-83-generic (hostname) 	01/01/2024 	_x86_64_	(4 CPU)
		
		Device             tps    kB_read/s    kB_wrtn/s    kB_dscd/s    kB_read    kB_wrtn    kB_dscd
		sda               9.68        12.34        34.56     123456     345678`,
			expected: 9.68,
			err:      false,
		},
		{
			name: "Localized numbers",
			input: `Linux 5.15.0-83-generic (hostname) 	01/01/2024 	_x86_64_	(4 CPU)

		Device             tps    kB_read/s    kB_wrtn/s    kB_dscd/s    kB_read    kB_wrtn    kB_dscd
		sda               9,68        12,34        34,56     123456     345678`,
			expected: 9.68,
			err:      false,
		},
		{
			name: "Multiple disks",
			input: `Linux 5.15.0-83-generic (hostname) 	01/01/2024 	_x86_64_	(4 CPU)

		Device             tps    kB_read/s    kB_wrtn/s    kB_dscd/s    kB_read    kB_wrtn    kB_dscd
		sda               1,00        12,34        34,56     123456     345678
		sdb               2,50        10,00        20,00     100000     200000`,
			expected: 3.5,
			err:      false,
		},
		{
			name:     "No data after header",
			input:    "",
			expected: 0,
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := parseDiskTransfersPerSecWithIostat(tt.input)
			if tt.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				fmt.Println("exp:", tt.expected, "res:", res)
				require.InEpsilon(t, tt.expected, res, 0.01)
			}
		})
	}
}

func TestParseDiskBytesPerSecWithIostat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
		err      bool
	}{
		{
			name: "One disk",
			input: `Linux 5.15.0-83-generic (hostname) 	01/01/2024 	_x86_64_	(4 CPU)

Device             tps    kB_read/s    kB_wrtn/s    kB_dscd/s    kB_read    kB_wrtn    kB_dscd
sda               9.68        53.16       294.52         0.00  188591750 1044819781          0`,
			expected: 347.68, // 53.16 + 294.52
			err:      false,
		},
		{
			name: "Localized numbers",
			input: `Linux 5.15.0-83-generic (hostname) 	01/01/2024 	_x86_64_	(4 CPU)

Device             tps    kB_read/s    kB_wrtn/s    kB_dscd/s    kB_read    kB_wrtn    kB_dscd
sda               9,68        53,16       294,52         0,00  188591750 1044819781          0`,
			expected: 347.68,
			err:      false,
		},
		{
			name: "Multiple disks",
			input: `Linux 5.15.0-83-generic (hostname) 	01/01/2024 	_x86_64_	(4 CPU)

Device             tps    kB_read/s    kB_wrtn/s    kB_dscd/s    kB_read    kB_wrtn    kB_dscd
sda               1,00        10,00        20,00         0,00   10000000    20000000          0
sdb               2,00        30,00        40,00         0,00   30000000    40000000          0`,
			expected: 100.00, // (10+20)+(30+40)
			err:      false,
		},
		{
			name:     "Empty input",
			input:    "",
			expected: 0,
			err:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := parseDiskBytesPerSecWithIostat(tt.input)
			if tt.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.InEpsilon(t, tt.expected, res, 0.01)
			}
		})
	}
}
