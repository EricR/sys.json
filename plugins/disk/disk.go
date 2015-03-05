package disk

import (
	"github.com/ericr/sysjson/util"
	"os"
	"syscall"
)

type DiskStats struct {
	Reads  DiskStatsForOp `json:"reads"`
	Writes DiskStatsForOp `json:"writes"`
	IO     DiskStatsForIO `json:"io"`
	Size   int            `json:"size"`
	Free   int            `json:"free"`
}
type DiskStatsForIO struct {
	InProgress      int `json:"in_progress"`
	TotalMs         int `json:"total_ms"`
	TotalWeightedMs int `json:"total_weighted"ms"`
}
type DiskStatsForOp struct {
	Completed int `json:"completed"`
	Sectors   int `json:"sectors"`
	Merged    int `json:"merged"`
	TotalMs   int `json:"total_ms"`
}

func GetDiskInfo() map[string]interface{} {
	stats := map[string]interface{}{}
	partitions := []string{}
	util.EachLine("/proc/partitions", func(fields []string) {
		if fields[0] != "major" {
			partitions = append(partitions, fields[3])
		}
	})
	util.EachLine("/proc/diskstats", func(fields []string) {
		for _, partition := range partitions {
			partSize, partFree, _ := Space(partition)
			if fields[2] == partition {
				stats[partition] = DiskStats{
					Reads: DiskStatsForOp{
						Completed: util.ParseInt(fields[3]),
						Merged:    util.ParseInt(fields[4]),
						Sectors:   util.ParseInt(fields[5]),
						TotalMs:   util.ParseInt(fields[6]),
					},
					Writes: DiskStatsForOp{
						Completed: util.ParseInt(fields[7]),
						Merged:    util.ParseInt(fields[8]),
						Sectors:   util.ParseInt(fields[9]),
						TotalMs:   util.ParseInt(fields[10]),
					},
					IO: DiskStatsForIO{
						InProgress:      util.ParseInt(fields[11]),
						TotalMs:         util.ParseInt(fields[11]),
						TotalWeightedMs: util.ParseInt(fields[12]),
					},
					Size: partSize,
					Free: partFree,
				}
			}
		}
	})
	return stats
}

func Space(path string) (total int, free int, err error) {
	var stat syscall.Statfs_t
	wd, err := os.Getwd()
	err = syscall.Statfs(wd, &stat)
	if err != nil {
		return
	}
	total = int(stat.Bsize) * int(stat.Blocks)
	free = int(stat.Bavail) * int(stat.Bsize)
	return
}

