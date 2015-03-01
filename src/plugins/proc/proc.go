package proc

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"unicode"
)

type FieldMap map[string]interface{}

func (fm FieldMap) Map(fn func(collector FieldMap, key string, value FieldMap)) map[string]interface{} {
	collector := map[string]interface{}{}

	for k, v := range fm {
		fn(collector, k, v.(FieldMap))
	}

	return collector
}

func GetLoadAvg() FieldMap {
	return getFieldMap(readFile("/proc/loadavg"), func(fm FieldMap, fields []string) {
		processes := strings.Split(fields[3], "/")

		fm["1m"] = parse(fields[0], "float")
		fm["5m"] = parse(fields[1], "float")
		fm["15m"] = parse(fields[2], "float")
		fm["prunning"] = parse(processes[0], "int")
		fm["lastpid"] = parse(fields[4], "int")
	})
}

func GetProcessesTree() FieldMap {
	fm := FieldMap{}
	pids := getPids()

	for _, pid := range pids {
		info := getProcessInfo(string(pid))
		fm[string(pid)] = FieldMap{
			"name":    info["Name"],
			"cmdline": info["CmdLine"],
			"state":   info["State"],
			"pid":     parse(info["Pid"].(string), "int"),
			"ppid":    parse(info["PPid"].(string), "int"),
			"threads": parse(info["Threads"].(string), "int"),
		}
	}

	return fm
}

func GetNetworkInfo() FieldMap {
	return getFieldMap(readFile("/proc/net/dev"), func(fm FieldMap, fields []string) {
		if !strings.Contains(fields[0], "|") && !strings.Contains(fields[1], "|") {
			if strings.Contains(fields[0], ":") {
				split := strings.Split(fields[0], ":")
				fields = append(split, fields[1:]...)
			}
			if fields[1] == "" {
				fields = append([]string{fields[0]}, fields[2:]...)
			}

			fm[strings.TrimRight(fields[0], ":")] = FieldMap{
				"receive_bytes":    parse(fields[1], "int"),
				"receive_packets":  parse(fields[2], "int"),
				"receive_errors":   parse(fields[3], "int"),
				"receive_drops":    parse(fields[4], "int"),
				"transmit_bytes":   parse(fields[9], "int"),
				"transmit_packets": parse(fields[10], "int"),
				"transmit_errors":  parse(fields[11], "int"),
				"transmit_drops":   parse(fields[12], "int"),
			}
		}
	})
}

func GetDiskInfo() FieldMap {
	var partitions []string

	getFieldMap(readFile("/proc/partitions"), func(fm FieldMap, fields []string) {
		if fields[0] != "major" {
			partitions = append(partitions, fields[3])
		}
	})

	return getFieldMap(readFile("/proc/diskstats"), func(fm FieldMap, fields []string) {
		for _, partition := range partitions {
			if fields[2] == partition {
				fm[partition] = FieldMap{
					"reads_completed":       parse(fields[3], "int"),
					"reads_merged":          parse(fields[4], "int"),
					"reads_sectors":         parse(fields[5], "int"),
					"reads_total_ms":        parse(fields[6], "int"),
					"writes_completed":      parse(fields[7], "int"),
					"writes_merged":         parse(fields[8], "int"),
					"writes_sectors":        parse(fields[9], "int"),
					"writes_total_ms":       parse(fields[10], "int"),
					"ios_in_progress":       parse(fields[11], "int"),
					"ios_total_ms":          parse(fields[11], "int"),
					"ios_total_weighted_ms": parse(fields[12], "int"),
				}
			}
		}
	})
}

func GetMemInfo() FieldMap {
	return getFieldMap(readFile("/proc/meminfo"), func(fm FieldMap, fields []string) {
		fm[strings.TrimRight(fields[0], ":")] = parse(fields[1], "int")
	})
}

func GetUptime() FieldMap {
	return getFieldMap(readFile("/proc/uptime"), func(fieldMap FieldMap, fields []string) {
		fieldMap["total"] = parse(fields[0], "float")
		fieldMap["idle"] = parse(fields[1], "float")
	})
}

// Private functions

func getPids() []string {
	pids := []string{}
	procFiles, _ := ioutil.ReadDir("/proc/")
	for _, procFile := range procFiles {
		if pid := parse(procFile.Name(), "int"); pid != 0 {
			pids = append(pids, procFile.Name())
		}
	}

	return pids
}

func getProcessInfo(pid string) FieldMap {
	fm := getFieldMap(readFile("/proc/"+pid+"/status"), func(fm FieldMap, fields []string) {
		if len(fields) > 1 {
			fm[strings.TrimRight(fields[0], ":")] = fields[1]
		}
	})

	cmdLine := strings.FieldsFunc(readFile("/proc/"+pid+"/cmdline"), func(r rune) bool {
		return !unicode.IsPrint(r)
	})
	fm["CmdLine"] = strings.Join(cmdLine, " ")

	return fm
}

func readFile(filePath string) string {
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("[error] Could not read file %s: %s", filePath, err)
		return ""
	}

	return string(contents)
}

func getFieldMap(s string, fn func(fm FieldMap, fields []string)) FieldMap {
	fieldMap := FieldMap{}

	for _, line := range strings.Split(s, "\n") {
		if len(line) > 0 {
			fn(fieldMap, strings.Fields(string(line)))
		}
	}
	return fieldMap
}

func parse(s, cType string) interface{} {
	switch cType {
	case "int":
		result, _ := strconv.Atoi(s)
		return result
	case "float":
		result, _ := strconv.ParseFloat(s, 64)
		return result
	}

	return string(s)
}
