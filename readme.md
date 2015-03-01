# sys.json

Exposes system stats as a JSON API.

It works, but it's a work in progress. Things are subject to change
without notice.

## Running It

```bash
./sysjson --listen 0.0.0.0:5374
```

## Todo

* Write tests
* HTTP authentication for security
* Standardize and document API
* Allow filtering so you only get back what you're interested in
* More in-depth networking and disk stats
* Plugin engine for reading, parsing, and presenting data
