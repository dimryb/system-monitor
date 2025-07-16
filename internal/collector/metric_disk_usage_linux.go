//go:build linux

package collector

import (
	"context"
	"fmt"

	"github.com/dimryb/system-monitor/internal/entity"
	i "github.com/dimryb/system-monitor/internal/interface"
)

type diskUsageMetric struct {
	value          *[]entity.DiskUsage
	collectorUsage i.ParamCollector
	parserUsage    func(rawData string) ([]entity.DiskUsage, error)

	collectorInode i.ParamCollector
	parserInode    func(rawData string) ([]entity.DiskUsage, error)
}

func (m *diskUsageMetric) collect(ctx context.Context) error {
	rawUsage, err := m.collectorUsage.Collect(ctx)
	if err != nil {
		return err
	}

	diskUsage, err := m.parserUsage(rawUsage)
	if err != nil {
		return err
	}

	rawInode, err := m.collectorInode.Collect(ctx)
	if err != nil {
		return err
	}

	diskInodes, err := m.parserInode(rawInode)
	if err != nil {
		return err
	}

	inodeMap := buildInodeMap(diskInodes)

	for i := range diskUsage {
		if inode, ok := inodeMap[diskUsage[i].Name]; ok {
			diskUsage[i].InodesTotal = inode.InodesTotal
			diskUsage[i].InodesUsed = inode.InodesUsed
			diskUsage[i].InodesUsedPercent = inode.InodesUsedPercent
		} else {
			diskUsage[i].InodesTotal = 0
			diskUsage[i].InodesUsed = 0
			diskUsage[i].InodesUsedPercent = -1
		}
	}

	*m.value = diskUsage

	fmt.Println(*m.value)
	return nil
}

func buildInodeMap(inodes []entity.DiskUsage) map[string]entity.DiskUsage {
	m := make(map[string]entity.DiskUsage)
	for _, disk := range inodes {
		m[disk.Name] = disk
	}
	return m
}
