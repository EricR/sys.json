package proc

import (
	"bytes"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"unicode"
)

type fieldMap map[string]interface{}

func GetLoadAvg() fieldMap {
	return getFieldMap(readFile("/proc/loadavg"), func(fieldMap fieldMap, fields []string) {
		processes := strings.Split(fields[3], "/")

		fieldMap["1m"] = parse(fields[0], "float")
		fieldMap["5m"] = parse(fields[1], "float")
		fieldMap["15m"] = parse(fields[2], "float")
		fieldMap["prunning"] = parse(processes[0], "int")
		fieldMap["ptotal"] = parse(processes[1], "int")
		fieldMap["lastpid"] = parse(fields[4], "int")
	})
}

func GetProcessesInfo() fieldMap {
	fm := fieldMap{}

	procFiles, _ := ioutil.ReadDir("/proc/")
	for _, procFile := range procFiles {
		if pid := parse(procFile.Name(), "int"); pid != 0 {
			cmdLine := string(readFile("/proc/" + procFile.Name() + "/cmdline"))
			strippedCmdLine := strings.FieldsFunc(cmdLine, func(r rune) bool {
				return !unicode.IsPrint(r)
			})
			status := getFieldMap(readFile("/proc/"+procFile.Name()+"/status"), func(fm fieldMap, fields []string) {
				if len(fields) > 1 {
					fm[strings.TrimRight(fields[0], ":")] = fields[1]
				}
			})
			fm[procFile.Name()] = fieldMap{
				"cmdline": strings.Join(strippedCmdLine, " "),
				"status":  status,
			}
		}
	}

	return fm
}

func GetDiskInfo() fieldMap {
	var partitions []string

	getFieldMap(readFile("/proc/partitions"), func(fm fieldMap, fields []string) {
		if fields[0] != "major" {
			partitions = append(partitions, fields[3])
		}
	})

	return getFieldMap(readFile("/proc/diskstats"), func(fm fieldMap, fields []string) {
		for _, partition := range partitions {
			if fields[2] == partition {
				fm[partition] = fieldMap{
					"reads": fieldMap{
						"completed": parse(fields[3], "int"),
						"merged":    parse(fields[4], "int"),
						"sectors":   parse(fields[5], "int"),
						"total_ms":  parse(fields[6], "int"),
					},
					"writes": fieldMap{
						"completed": parse(fields[7], "int"),
						"merged":    parse(fields[8], "int"),
						"sectors":   parse(fields[9], "int"),
						"total_ms":  parse(fields[10], "int"),
					},
					"ios_in_progress":       parse(fields[11], "int"),
					"ios_total_ms":          parse(fields[11], "int"),
					"ios_total_weighted_ms": parse(fields[12], "int"),
				}
			}
		}
	})
}

func GetMemInfo() fieldMap {
	return getFieldMap(readFile("/proc/meminfo"), func(fm fieldMap, fields []string) {
		fm[strings.TrimRight(fields[0], ":")] = parse(fields[1], "int")
	})
}

func GetUptime() fieldMap {
	return getFieldMap(readFile("/proc/uptime"), func(fieldMap fieldMap, fields []string) {
		fieldMap["total"] = parse(fields[0], "float")
		fieldMap["idle"] = parse(fields[1], "float")
	})
}

// Private functions

func readFile(filePath string) []byte {
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("[error] Could not read file %s: %s", filePath, err)
		return []byte{}
	}

	return contents
}

func getFieldMap(b []byte, fn func(fieldMap fieldMap, fields []string)) fieldMap {
	fieldMap := fieldMap{}

	for _, line := range bytes.Split(b, []byte("\n")) {
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
		result, _ := strconv.ParseFloat(string(s), 64)
		return result
	}

	return string(s)
}
