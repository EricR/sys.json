package host

import (
	"os"
	"time"
)

func GetInfo() map[string]interface{} {
	hostname, _ := os.Hostname()
	stime := time.Now()

	return map[string]interface{}{
		"name": hostname,
		"time": map[string]interface{}{
			"string": stime,
			"unix":   stime.Unix(),
		},
	}
}
