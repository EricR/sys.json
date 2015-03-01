# sys.json

Exposes system stats as a JSON API.

It works, but it's a work in progress. Things are subject to change
without notice.

## Why?

There are plenty of other tools out there that can give you stats. There's SNMP, dozens of monitoring
solutions and stat collection agents, and more.

All of these solutions vary drastically in complexity and ease of use, many times sacrificing
one for the other. The goal of sys.json is to provide a simple API for querying a server about the
stats you're interested in. No more, no less. What you do with that data is up to you. Some ideas include:

* Alerts
* Creating graphs
* Capacity analysis

## Running It

```bash
./sysjson --listen 0.0.0.0:5374
```

## Example Response

```js
{
  "disk":{
    "vda":{
      "io_ops":{
        "current":0,
        "total_ms":0,
        "weighted_total_ms":25909
      },
      "reads":{
        "completed":9336,
        "merged":4151,
        "sectors":616650,
        "total_ms":2844
      },
      "writes":{
        "completed":39830,
        "merged":188632,
        "sectors":1828512,
        "total_ms":86448
      }
    },
    "vda1":{
      "io_ops":{
        "current":0,
        "total_ms":0,
        "weighted_total_ms":25903
      },
      "reads":{
        "completed":9162,
        "merged":4151,
        "sectors":615258,
        "total_ms":2837
      },
      "writes":{
        "completed":39830,
        "merged":188632,
        "sectors":1828512,
        "total_ms":86448
      }
    }
  },
  "load_avg":{
    "15m":0,
    "1m":0.07,
    "5m":0.03
  },
  "memory":{
    "active":288240,
    "active_anon":50084,
    "active_file":238156,
    "anon_huge_pages":20480,
    "anon_pages":50124,
    "bounce":0,
    "commit_limit":510200,
    "commited_as":270968,
    "direct_map":{
      "1G":0,
      "2M":1042432,
      "4k":6136
    },
    "dirty":20,
    "huge_page_size":2048,
    "huge_pages":{
      "free":0,
      "rsvd":0,
      "surp":0,
      "total":0
    },
    "inactive":165572,
    "inactive_anon":148,
    "inactive_file":165424,
    "kernel_stack":712,
    "mapped":47440,
    "mlocked":0,
    "nfs_unstable":0,
    "s_reclaimable":34752,
    "s_unreclaim":19244,
    "shmem":160,
    "simple":{
      "buffers":77664,
      "cached":326080,
      "free":495200,
      "total":1020400
      "swap":{
        "free":0,
        "total":0,
        "cached":0
      },
      "free_total":898944
    },
    "slab":53996,
    "swap":{
      "free":0,
      "total":0
    },
    "unevictable":0,
    "vmalloc":{
      "chunk":34359727808,
      "total":34359738367,
      "used":7052
    },
    "writeback":0,
    "writeback_tmp":0
  },
  "network":{
    "eth0":{
      "receive":{
        "bytes":20471337,
        "dropped":0,
        "errors":0,
        "packets":157255
      },
      "transmit":{
        "bytes":128487885,
        "dropped":0,
        "errors":0,
        "packets":134811
      }
    },
    "eth1":{
      "receive":{
        "bytes":0,
        "dropped":0,
        "errors":0,
        "packets":0
      },
      "transmit":{
        "bytes":566,
        "dropped":0,
        "errors":0,
        "packets":9
      }
    }
  },
  "processes":{
    "last_pid":28604,
    "running":2,
    "tree":{
      "1":{
        "cmdline":"/sbin/init",
        "name":"init",
        "pid":1,
        "ppid":0,
        "state":"S",
        "thread_count":1
      },
      "2":{
        "cmdline":"",
        "name":"kthreadd",
        "pid":2,
        "ppid":0,
        "state":"S",
        "thread_count":1
      },
      "3":{
        "cmdline":"",
        "name":"migration/0",
        "pid":3,
        "ppid":2,
        "state":"S",
        "thread_count":1
      },
      "4":{
        "cmdline":"",
        "name":"ksoftirqd/0",
        "pid":4,
        "ppid":2,
        "state":"S",
        "thread_count":1
      },
      "5":{
        "cmdline":"",
        "name":"migration/0",
        "pid":5,
        "ppid":2,
        "state":"S",
        "thread_count":1
      },
      "6":{
        "cmdline":"",
        "name":"watchdog/0",
        "pid":6,
        "ppid":2,
        "state":"S",
        "thread_count":1
      }
    }
  },
  "uptime":{
    "idle":74001.81,
    "total":75057.54
  }
}
```

## Todo

* Write tests
* HTTP authentication for security
* Standardize and document API
* Allow filtering so you only get back what you're interested in
* More in-depth networking and disk stats
* Plugin engine for reading, parsing, and presenting data
