package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

var (
	listenAddr = flag.String("listen", "0.0.0.0:6535", "Address to listen on.")
)

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", statsHandler)

	log.Printf("[notice] Starting sys.json on %s...", *listenAddr)
	http.ListenAndServe(*listenAddr, mux)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	loadAvgs := getLoadAvgs()
	memory := getMemory()
	uptime := getUptime()

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
			"free_total":  memory["MemFree"].(int) + memory["Buffers"].(int) + memory["Cached"].(int),
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

func getLoadAvgs() map[string]interface{} {
	return fileNamedMatchesWithTypeConversions("/proc/loadavg",
		`(?P<1m>\d+\.\d+) (?P<5m>\d+\.\d+) (?P<15m>\d+\.\d+) (?P<prunning>\d+)\/(?P<ptotal>\d+) (?P<lastpid>\d+)`,
		map[string]string{"1m": "float", "5m": "float", "15m": "float", "prunning": "int", "ptotal": "int", "lastpid": "int"})
}

func getMemory() map[string]interface{} {
	return fileKeyValueMatches("/proc/meminfo")
}

func getUptime() map[string]interface{} {
	return fileNamedMatchesWithTypeConversions("/proc/uptime",
		`(?P<total>\d+\.\d+) (?P<idle>\d+\.\d+)`,
		map[string]string{"total": "float", "idle": "float"})
}

func fileKeyValueMatches(file string) map[string]interface{} {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("[error] Could not read file %s: %s", file, err)
		return map[string]interface{}{}
	}

	regex := `(\S+):[\s]+(\d+)`
	re, err2 := regexp.Compile(regex)
	if err2 != nil {
		log.Printf("[error] Could not compile regexp: %s", err2)
		return map[string]interface{}{}
	}

	submatches := re.FindAllSubmatch(content, -1)
	matches := map[string]interface{}{}

	for _, n := range submatches {
		r, _ := strconv.Atoi(string(n[2]))
		matches[string(n[1])] = r
	}

	return matches
}

func fileNamedMatchesWithTypeConversions(file, regex string, convMap map[string]string) map[string]interface{} {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("[error] Could not read file: %s: %s", file, err)
		return map[string]interface{}{}
	}

	re, err2 := regexp.Compile(regex)
	if err2 != nil {
		log.Printf("[error] Could not compile regexp: %s", err2)
		return map[string]interface{}{}
	}

	names := re.SubexpNames()
	submatches := re.FindAllSubmatch(content, -1)[0]
	matches := map[string]interface{}{}

	for i, n := range submatches {
		if i != 0 {
			name := names[i]
			if len(convMap) > 0 {
				switch convMap[name] {
				case "string":
					matches[name] = string(n)
				case "int":
					r, _ := strconv.Atoi(string(n))
					matches[name] = r
				case "float":
					r, _ := strconv.ParseFloat(string(n), 64)
					matches[name] = r
				}
			} else {
				matches[name] = string(n)
			}
		}
	}

	return matches
}
