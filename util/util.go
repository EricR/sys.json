package util

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func ParseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func ParseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func EachLine(file string, fn func([]string)) {
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("[error] Could not read file: %s: %s", file, err)
		return
	}

	for _, line := range strings.Split(string(contents), "\n") {
		fields := strings.Fields(line)
		if len(fields) > 1 {
			fn(fields)
		}
	}
}

func ReadFile(file string) string {
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("[error] Could not read file: %s: %s", file, err)
	}

	return string(contents)
}
