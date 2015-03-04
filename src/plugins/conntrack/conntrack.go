package conntrack

import (
	"bufio"
	"github.com/ericr/sysjson/src/util"
	"log"
	"os/exec"
	"strings"
)

type Connection struct {
	State       string         `json:"state"`
	Source      ConnectionSide `json:"source"`
	Destination ConnectionSide `json:"destination"`
}

type ConnectionSide struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

func GetStats() map[string]interface{} {
	stats := map[string]interface{}{}
	connections := []Connection{}
	counters := map[string]int{}

	path, err := exec.LookPath("conntrack")
	if err != nil {
		log.Printf("[error] conntrack: %s", err)
		return stats
	}

	cmd := exec.Command(path, "-L -f -o extend")

	out, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("[error] conntrack: %", err)
		return stats
	}

	if err := cmd.Start(); err != nil {
		log.Printf("[error] conntrack: %s", err)
		return stats
	}

	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())

			state := strings.ToLower(fields[3])
			counters[state] = counters[state] + 1

			src := strings.Split(fields[4], "=")[1]
			dst := strings.Split(fields[5], "=")[1]
			sport := strings.Split(fields[6], "=")[1]
			dport := strings.Split(fields[7], "=")[1]

			conn := Connection{
				State: state,
				Source: ConnectionSide{
					Address: src,
					Port:    util.ParseInt(sport),
				},
				Destination: ConnectionSide{
					Address: dst,
					Port:    util.ParseInt(dport),
				},
			}

			connections = append(connections, conn)
		}
		if err := scanner.Err(); err != nil {
			log.Printf("[error] conntrack: %s", err)
			return stats
		}
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("[error] conntrack: %s", err)
		return stats
	}

	stats["connections"] = connections
	stats["counters"] = counters

	return stats
}
