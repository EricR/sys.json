package main

import (
	"encoding/json"
	"github.com/ericr/sysjson/src/plugins/proc"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", statsHandler)
	http.ListenAndServe(":3000", mux)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	loadAvgs := proc.GetLoadAvg()
	memory := proc.GetMemInfo()
	uptime := proc.GetUptime()

	resp := map[string]interface{}{}
	resp["uptime"] = map[string]interface{}{
		"total": uptime["total"],
		"idle":  uptime["idle"],
	}
	resp["load_avgs"] = map[string]interface{}{
		"1m":  loadAvgs["1m"],
		"5m":  loadAvgs["5m"],
		"15m": loadAvgs["15m"],
	}
	resp["processes"] = map[string]interface{}{
		"running":  loadAvgs["prunning"],
		"sleeping": loadAvgs["ptotal"].(int) - loadAvgs["prunning"].(int),
		"total":    loadAvgs["ptotal"],
		"last_pid": loadAvgs["lastpid"],
	}
	resp["memory"] = map[string]interface{}{
		"simple": map[string]interface{}{
			"total":       memory["MemTotal"],
			"free":        memory["MemFree"],
			"buffers":     memory["Buffers"],
			"cached":      memory["Cached"],
			"swap_cached": memory["SwapCached"],
			//"free_total":  memory["MemFree"].(int) + memory["Buffers"].(int) + memory["Cached"].(int),
		},
		"active":        memory["Active"],
		"inactive":      memory["Inactive"],
		"active_anon":   memory["Active(anon)"],
		"inactive_anon": memory["Inactive(anon)"],
		"active_file":   memory["Active(file)"],
		"inactive_file": memory["Inactive(file)"],
		"unevictable":   memory["Unevictable"],
		"mlocked":       memory["Mlocked"],
		"swap": map[string]interface{}{
			"total": memory["SwapTotal"],
			"free":  memory["SwapFree"],
		},
		"dirty":         memory["Dirty"],
		"writeback":     memory["Writeback"],
		"anon_pages":    memory["AnonPages"],
		"mapped":        memory["Mapped"],
		"shmem":         memory["Shmem"],
		"slab":          memory["Slab"],
		"s_reclaimable": memory["SReclaimable"],
		"s_unreclaim":   memory["SUnreclaim"],
		"kernel_stack":  memory["KernelStack"],
		"nfs_unstable":  memory["NFS_Unstable"],
		"bounce":        memory["Bounce"],
		"writeback_tmp": memory["WritebackTmp"],
		"commit_limit":  memory["CommitLimit"],
		"commited_as":   memory["Committed_AS"],
		"vmalloc": map[string]interface{}{
			"total": memory["VmallocTotal"],
			"used":  memory["VmallocUsed"],
			"chunk": memory["VmallocChunk"],
		},
		"anon_huge_pages": memory["AnonHugePages"],
		"huge_pages": map[string]interface{}{
			"total": memory["HugePages_Total"],
			"free":  memory["HugePages_Free"],
			"rsvd":  memory["HugePages_Rsvd"],
			"surp":  memory["HugePages_Surp"],
		},
		"huge_page_size": memory["Hugepagesize"],
		"direct_map": map[string]interface{}{
			"4k": memory["DirectMap4k"],
			"2M": memory["DirectMap2M"],
			"1G": memory["DirectMap1G"],
		},
	}

	respJSON, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}
