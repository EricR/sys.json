package main

import (
	"github.com/ericr/sysjson/src/plugins/proc"
	"log"
	"strings"
)

type module struct {
	Name string
	Desc string
	Func func() j
}

type j map[string]interface{}

var (
	modules = map[string]module{
		"load": module{
			Name: "load",
			Desc: "Provides load averages",
			Func: func() j {
				load := proc.GetLoadAvg()
				return j{
					"1m":  load["1m"],
					"5m":  load["5m"],
					"15m": load["15m"],
				}
			},
		},
		"uptime": module{
			Name: "uptime",
			Desc: "Provides uptime",
			Func: func() j {
				uptime := proc.GetUptime()
				return j{
					"total": uptime["total"],
					"idle":  uptime["idle"],
				}
			},
		},
		"process": module{
			Name: "process",
			Desc: "Provides a process tree",
			Func: func() j {
				return proc.GetProcessTree().Map(func(fm proc.FieldMap, key string, val proc.FieldMap) {
					fm[key] = j{
						"name":         val["name"],
						"cmdline":      val["cmdline"],
						"state":        val["state"],
						"pid":          val["pid"],
						"ppid":         val["ppid"],
						"thread_count": val["threads"],
						"mem": j{
							"virtual": j{
								"current": val["vm_size"],
								"peak":    val["vm_peak"],
							},
							"resident": j{
								"current": val["vm_rss"],
								"peak":    val["vm_hwm"],
							},
						},
					}
				})
			},
		},
		"mem": module{
			Name: "mem",
			Desc: "Provides memory stats",
			Func: func() j {
				mem := proc.GetMemInfo()
				return j{
					"simple": j{
						"total":   mem["MemTotal"],
						"free":    mem["MemFree"],
						"buffers": mem["Buffers"],
						"cached":  mem["Cached"],
						"swap": j{
							"cached": mem["SwapCached"],
							"total":  mem["SwapTotal"],
							"free":   mem["SwapFree"],
						},
						"free_total": mem["MemFree"].(int) + mem["Buffers"].(int) + mem["Cached"].(int),
					},
					"active":        mem["Active"],
					"inactive":      mem["Inactive"],
					"active_anon":   mem["Active(anon)"],
					"inactive_anon": mem["Inactive(anon)"],
					"active_file":   mem["Active(file)"],
					"inactive_file": mem["Inactive(file)"],
					"unevictable":   mem["Unevictable"],
					"mlocked":       mem["Mlocked"],
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
			},
		},
		"disk": module{
			Name: "disk",
			Desc: "Provides stats on each disk",
			Func: func() j {
				disk := proc.GetDiskInfo()
				return disk.Map(func(fm proc.FieldMap, key string, val proc.FieldMap) {
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
			},
		},
		"net": module{
			Name: "net",
			Desc: "Provides stats on each network interface",
			Func: func() j {
				net := proc.GetNetworkInfo()
				return net.Map(func(fm proc.FieldMap, key string, val proc.FieldMap) {
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
			},
		},
	}
)

func loadModules(resp j, paramModules string) j {
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
		resp["modules_available"] = modulesAvailable
	}

	for _, m := range modulesSelected {
		module := modules[m]
		if module.Name != "" {
			resp[m] = module.Func()
		} else {
			log.Printf("[warn] Module '%s' not found", m)
		}
	}

	return resp
}
