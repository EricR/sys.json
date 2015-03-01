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

type j map[string]interface{}

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
	network := proc.GetNetworkInfo()

	resp := j{}
	resp["uptime"] = j{
		"total": uptime["total"],
		"idle":  uptime["idle"],
	}
	resp["load_avg"] = j{
		"1m":  load["1m"],
		"5m":  load["5m"],
		"15m": load["15m"],
	}
	resp["processes"] = j{
		"tree": processTree.Map(func(fm proc.FieldMap, key string, val proc.FieldMap) {
			fm[key] = j{
				"name":         val["name"],
				"cmdline":      val["cmdline"],
				"state":        val["state"],
				"pid":          val["pid"],
				"ppid":         val["ppid"],
				"thread_count": val["threads"],
			}
		}),
		"running":  load["prunning"],
		"last_pid": load["lastpid"],
	}
	resp["memory"] = j{
		"simple": j{
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
		"swap": j{
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
		"vmalloc": j{
			"total": mem["VmallocTotal"],
			"used":  mem["VmallocUsed"],
			"chunk": mem["VmallocChunk"],
		},
		"anon_huge_pages": mem["AnonHugePages"],
		"huge_pages": j{
			"total": mem["HugePages_Total"],
			"free":  mem["HugePages_Free"],
			"rsvd":  mem["HugePages_Rsvd"],
			"surp":  mem["HugePages_Surp"],
		},
		"huge_page_size": mem["Hugepagesize"],
		"direct_map": j{
			"4k": mem["DirectMap4k"],
			"2M": mem["DirectMap2M"],
			"1G": mem["DirectMap1G"],
		},
	}
	resp["disk"] = disk.Map(func(fm proc.FieldMap, key string, val proc.FieldMap) {
		fm[key] = j{
			"reads": j{
				"completed": val["reads_completed"],
				"merged":    val["reads_merged"],
				"sectors":   val["reads_sectors"],
				"total_ms":  val["reads_total_ms"],
			},
			"writes": j{
				"completed": val["writes_completed"],
				"merged":    val["writes_merged"],
				"sectors":   val["writes_sectors"],
				"total_ms":  val["writes_total_ms"],
			},
			"io_ops": j{
				"current":           val["ios_in_progress"],
				"total_ms":          val["ios_total_ms"],
				"weighted_total_ms": val["ios_total_weighted_ms"],
			},
		}
	})
	resp["network"] = network.Map(func(fm proc.FieldMap, key string, val proc.FieldMap) {
		fm[key] = j{
			"receive": j{
				"bytes":   val["receive_bytes"],
				"packets": val["receive_packets"],
				"errors":  val["receive_errors"],
				"dropped": val["receive_drops"],
			},
			"transmit": j{
				"bytes":   val["transmit_bytes"],
				"packets": val["transmit_packets"],
				"errors":  val["transmit_errors"],
				"dropped": val["transmit_drops"],
			},
		}
	})

	respJSON, err := json.Marshal(resp)
	if err != nil {
		log.Fatal("[error] Fatal! Could not construct JSON response: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}
