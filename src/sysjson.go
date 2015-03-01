package main

import (
	"encoding/json"
	"flag"
	"github.com/ericr/sysjson/src/plugins/proc"
	"log"
	"net/http"
)

var (
	listen = flag.String("listen", ":5374", "Address to listen on.")
)

func main() {
	flag.Parse()
	log.Printf("[notice] sys.json listening on %s", *listen)

	mux := http.NewServeMux()
	mux.HandleFunc("/", statsHandler)
	http.ListenAndServe(*listen, mux)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	load := proc.GetLoadAvg()
	mem := proc.GetMemInfo()
	uptime := proc.GetUptime()
	disk := proc.GetDiskInfo()
	processTree := proc.GetProcessesTree()

	resp := map[string]interface{}{}
	resp["uptime"] = map[string]interface{}{
		"total": uptime["total"],
		"idle":  uptime["idle"],
	}
	resp["load_avg"] = map[string]interface{}{
		"1m":  load["1m"],
		"5m":  load["5m"],
		"15m": load["15m"],
	}
	resp["processes"] = map[string]interface{}{
		"tree":     processTree,
		"running":  load["prunning"],
		"sleeping": load["ptotal"].(int) - load["prunning"].(int),
		"total":    load["ptotal"],
		"last_pid": load["lastpid"],
	}
	resp["memory"] = map[string]interface{}{
		"simple": map[string]interface{}{
			"total":       mem["MemTotal"],
			"free":        mem["MemFree"],
			"buffers":     mem["Buffers"],
			"cached":      mem["Cached"],
			"swap_cached": mem["SwapCached"],
			"free_total":  mem["MemFree"].(int) + mem["Buffers"].(int) + mem["Cached"].(int),
		},
		"active":        mem["Active"],
		"inactive":      mem["Inactive"],
		"active_anon":   mem["Active(anon)"],
		"inactive_anon": mem["Inactive(anon)"],
		"active_file":   mem["Active(file)"],
		"inactive_file": mem["Inactive(file)"],
		"unevictable":   mem["Unevictable"],
		"mlocked":       mem["Mlocked"],
		"swap": map[string]interface{}{
			"total": mem["SwapTotal"],
			"free":  mem["SwapFree"],
		},
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
	resp["disk"] = disk

	respJSON, err := json.Marshal(resp)
	if err != nil {
		log.Fatal("[error] Fatal! Could not construct JSON response: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}
