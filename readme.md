# sys.json

Exposes system stats as a JSON API.

It works, but it's a work in progress. Things are subject to change without notice, so please use
caution if you use this in production. Assume there are bugs.

## Why?

There are plenty of other tools out there that can give you stats. There's SNMP, dozens of monitoring
solutions and stat collection agents, and more. So why this?

Many of these solutions are complicated to setup and get data from. The goal of sys.json is to
provide a simple API for querying a server about the stats you're interested in. No more, no less.
What you do with that data is up to you.

## Running It

```bash
./sysjson --listen 0.0.0.0:5374
```

## Example Response

```bash
curl http://localhost:5374/?modules=disk,load
```

```js
{
  "hostname":"eric.localhost",
  "current_time":{
    "string":"2015-03-01T23:45:45.286141452-05:00",
    "unix":1425271545
  },
  "disk":{
    "vda":{
      "io_ops":{
        "current":0,
        "total_ms":0,
        "weighted_total_ms":34522
      },
      "reads":{
        "completed":9347,
        "merged":4151,
        "sectors":617802,
        "total_ms":2849
      },
      "writes":{
        "completed":53196,
        "merged":233210,
        "sectors":2292064,
        "total_ms":134061
      }
    },
    "vda1":{
      "io_ops":{
        "current":0,
        "total_ms":0,
        "weighted_total_ms":34516
      },
      "reads":{
        "completed":9173,
        "merged":4151,
        "sectors":616410,
        "total_ms":2842
      },
      "writes":{
        "completed":53196,
        "merged":233210,
        "sectors":2292064,
        "total_ms":134061
      }
    }
  },
  "load":{
    "15m":0,
    "1m":0.02,
    "5m":0.04
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
