package proc

import (
	"github.com/ericr/sysjson/src/util"
	"io/ioutil"
	"strings"
	"unicode"
)

type DiskStats struct {
	Reads  DiskStatsForOp `json:"reads"`
	Writes DiskStatsForOp `json:"writes"`
	IO     DiskStatsForIO `json:"io"`
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

type NetworkStats struct {
	Receieve NetworkStatsSide `json:"receieve"`
	Send     NetworkStatsSide `json:"send"`
}

type NetworkStatsSide struct {
	Bytes   int `json:"bytes"`
	Packets int `json:"packets"`
	Errors  int `json:"errors"`
	Dropped int `json:"dropped"`
}

type Process struct {
	Name        string        `json:"name"`
	CmdLine     string        `json:"cmdline"`
	State       string        `json:"state"`
	PID         int           `json:"pid"`
	PPID        int           `json:"ppid"`
	ThreadCount int           `json:"thread_count"`
	Memory      ProcessMemory `json:"memory"`
}

type ProcessMemory struct {
	Virtual  ProcessMemoryForType `json:"virtual"`
	Resident ProcessMemoryForType `json:"resident"`
}

type ProcessMemoryForType struct {
	Current int `json: "current"`
	Peak    int `json:"peak"`
}

func GetLoadAvg() map[string]interface{} {
	stats := map[string]interface{}{}

	util.EachLine("/proc/loadavg", func(fields []string) {
		stats["1m"] = util.ParseFloat(fields[0])
		stats["5m"] = util.ParseFloat(fields[1])
		stats["15m"] = util.ParseFloat(fields[2])
	})

	return stats
}

func GetProcessTree() map[string]interface{} {
	tree := map[string]interface{}{}
	pids := getPids()

	for _, pid := range pids {
		pids := string(pid)
		info := getProcessInfo(pids)
		tree[pids] = Process{
			Name:        info["Name"],
			CmdLine:     info["CmdLine"],
			State:       info["State"],
			PID:         util.ParseInt(info["Pid"]),
			PPID:        util.ParseInt(info["PPid"]),
			ThreadCount: util.ParseInt(info["Threads"]),
			Memory: ProcessMemory{
				Virtual: ProcessMemoryForType{
					Current: util.ParseInt(info["VmSize"]),
					Peak:    util.ParseInt(info["VmPeak"]),
				},
				Resident: ProcessMemoryForType{
					Current: util.ParseInt(info["VmRSS"]),
					Peak:    util.ParseInt(info["VmHWM"]),
				},
			},
		}
	}

	return tree
}

func GetNetworkInfo() map[string]interface{} {
	stats := map[string]interface{}{}

	util.EachLine("/proc/net/dev", func(fields []string) {
		if !strings.Contains(fields[0], "|") && !strings.Contains(fields[1], "|") {
			if strings.Contains(fields[0], ":") {
				split := strings.Split(fields[0], ":")
				fields = append(split, fields[1:]...)
			}
			if fields[1] == "" {
				fields = append([]string{fields[0]}, fields[2:]...)
			}

			stats[strings.TrimRight(fields[0], ":")] = NetworkStats{
				Receieve: NetworkStatsSide{
					Bytes:   util.ParseInt(fields[1]),
					Packets: util.ParseInt(fields[2]),
					Errors:  util.ParseInt(fields[3]),
					Dropped: util.ParseInt(fields[4]),
				},
				Send: NetworkStatsSide{
					Bytes:   util.ParseInt(fields[9]),
					Packets: util.ParseInt(fields[10]),
					Errors:  util.ParseInt(fields[11]),
					Dropped: util.ParseInt(fields[12]),
				},
			}
		}
	})

	return stats
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
				}
			}
		}
	})

	return stats
}

func GetMemoryInfo() map[string]interface{} {
	mem := map[string]int{}

	util.EachLine("/proc/meminfo", func(fields []string) {
		mem[strings.TrimRight(fields[0], ":")] = util.ParseInt(fields[1])
	})

	return map[string]interface{}{
		"simple": map[string]interface{}{
			"total":   mem["MemTotal"],
			"free":    mem["MemFree"],
			"buffers": mem["Buffers"],
			"cached":  mem["Cached"],
			"swap": map[string]int{
				"cached": mem["SwapCached"],
				"total":  mem["SwapTotal"],
				"free":   mem["SwapFree"],
			},
			"free_total": mem["MemFree"] + mem["Buffers"] + mem["Cached"],
		},
		"active":        mem["Active"],
		"inactive":      mem["Inactive"],
		"active_anon":   mem["Active(anon)"],
		"inactive_anon": mem["Inactive(anon)"],
		"active_file":   mem["Active(file)"],
		"inactive_file": mem["Inactive(file)"],
		"unevictable":   mem["Unevictable"],
		"mlocked":       mem["Mlocked"],
		"dirty":         mem["Dirty"],
		"writeback":     mem["Writeback"],
		"anon_pages":    mem["AnonPages"],
		"mapped":        mem["Mapped"],
		"shmem":         mem["Shmem"],
		"slab":          mem["Slab"],
		"s_reclaimable": mem["SReclaimable"],
		"s_unreclaim":   mem["SUnreclaim"],
		"kernel_stack":  mem["KernelStack"],
		"nfs_unstable":  mem["NFS_Unstable"],
		"bounce":        mem["Bounce"],
		"writeback_tmp": mem["WritebackTmp"],
		"commit_limit":  mem["CommitLimit"],
		"commited_as":   mem["Committed_AS"],
		"vmalloc": map[string]interface{}{
			"total": mem["VmallocTotal"],
			"used":  mem["VmallocUsed"],
			"chunk": mem["VmallocChunk"],
		},
		"anon_huge_pages": mem["AnonHugePages"],
		"huge_pages": map[string]interface{}{
			"total": mem["HugePages_Total"],
			"free":  mem["HugePages_Free"],
			"rsvd":  mem["HugePages_Rsvd"],
			"surp":  mem["HugePages_Surp"],
		},
		"huge_page_size": mem["Hugepagesize"],
		"direct_map": map[string]interface{}{
			"4k": mem["DirectMap4k"],
			"2M": mem["DirectMap2M"],
			"1G": mem["DirectMap1G"],
		},
	}
}

func GetUptime() map[string]interface{} {
	stats := map[string]interface{}{}

	util.EachLine("/proc/uptime", func(fields []string) {
		stats["total"] = util.ParseFloat(fields[0])
		stats["idle"] = util.ParseFloat(fields[1])
	})

	return stats
}

func getPids() []string {
	pids := []string{}
	procFiles, _ := ioutil.ReadDir("/proc/")

	for _, procFile := range procFiles {
		if pid := util.ParseInt(procFile.Name()); pid != 0 {
			pids = append(pids, procFile.Name())
		}
	}

	return pids
}

func getProcessInfo(pid string) map[string]string {
	info := map[string]string{}

	util.EachLine("/proc/"+pid+"/status", func(fields []string) {
		info[strings.TrimRight(fields[0], ":")] = fields[1]
	})

	cmdLine := strings.FieldsFunc(util.ReadFile("/proc/"+pid+"/cmdline"), func(r rune) bool {
		return !unicode.IsPrint(r)
	})
	info["CmdLine"] = strings.Join(cmdLine, " ")

	return info
}
