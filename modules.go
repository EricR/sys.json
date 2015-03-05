package main

import (
	"github.com/ericr/sysjson/plugins/conntrack"
	"github.com/ericr/sysjson/plugins/host"
	"github.com/ericr/sysjson/plugins/proc"
	"log"
	"strings"
)

type module struct {
	Name string `json:"name"`
	Desc string `json:"description"`
	fn   func() map[string]interface{}
}

var (
	modules = map[string]module{
		"host": module{
			Name: "host",
			Desc: "Provides basic host info",
			fn:   host.GetInfo,
		},
		"load": module{
			Name: "load",
			Desc: "Provides load averages",
			fn:   proc.GetLoadAvg,
		},
		"uptime": module{
			Name: "uptime",
			Desc: "Provides time since startup",
			fn:   proc.GetUptime,
		},
		"process": module{
			Name: "process",
			Desc: "Provides a process tree",
			fn:   proc.GetProcessTree,
		},
		"mem": module{
			Name: "mem",
			Desc: "Provides system-wide memory stats",
			fn:   proc.GetMemoryInfo,
		},
		"disk": module{
			Name: "disk",
			Desc: "Provides stats on each disk",
			fn:   proc.GetDiskInfo,
		},
		"net": module{
			Name: "net",
			Desc: "Provides stats on each network interface",
			fn:   proc.GetNetworkInfo,
		},
		"conntrack": module{
			Name: "conntrack",
			Desc: "Provides stats on current network connections (requires conntrack-tools)",
			fn:   conntrack.GetStats,
		},
	}
)

func loadModules(resp map[string]interface{}, paramModules string) map[string]interface{} {
	modulesSelected := []string{}
	modulesAvailable := make([]string, 0, len(modules))

	for module := range modules {
		modulesAvailable = append(modulesAvailable, module)
	}

	if paramModules == "all" {
		modulesSelected = modulesAvailable
	} else if len(paramModules) > 0 {
		modulesSelected = strings.Split(paramModules, ",")
	} else {
		resp["modules_available"] = modules
	}

	for _, m := range modulesSelected {
		module := modules[m]
		if module.Name != "" {
			resp[m] = module.fn()
		} else {
			log.Printf("[warn] Module '%s' not found", m)
		}
	}

	return resp
}
