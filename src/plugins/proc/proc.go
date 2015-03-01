package proc

import (
	"bytes"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type fieldMap map[string]interface{}

func GetLoadAvg() fieldMap {
	return getFieldMap(readFile("/proc/loadavg"), func(fieldMap fieldMap, fields []string) {
		processes := strings.Split(fields[3], "/")

		fieldMap["1m"] = conv(fields[0], "float")
		fieldMap["5m"] = conv(fields[1], "float")
		fieldMap["15m"] = conv(fields[2], "float")
		fieldMap["prunning"] = conv(processes[0], "int")
		fieldMap["ptotal"] = conv(processes[1], "int")
		fieldMap["lastpid"] = conv(fields[4], "int")
	})
}

func GetMemInfo() fieldMap {
	return getFieldMap(readFile("/proc/meminfo"), func(fieldMap fieldMap, fields []string) {
		fieldMap[strings.TrimRight(fields[0], ":")] = conv(fields[1], "int")
	})
}

func GetUptime() fieldMap {
	return getFieldMap(readFile("/proc/uptime"), func(fieldMap fieldMap, fields []string) {
		fieldMap["total"] = conv(fields[0], "float")
		fieldMap["idle"] = conv(fields[1], "float")
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

func conv(s, cType string) interface{} {
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
