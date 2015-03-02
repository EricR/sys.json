package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	listen = flag.String("listen", ":5374", "Address to listen on")
	tls    = flag.Bool("tls", false, "Use TLS (requires -cert and -key)")
	cert   = flag.String("cert", "", "TLS cert file")
	key    = flag.String("key", "", "TLS key file")
)

func main() {
	flag.Parse()
	log.Printf("[notice] sys.json listening on %s", *listen)

	mux := http.NewServeMux()
	mux.HandleFunc("/", statsHandler)
	if *tls {
		log.Printf("[notice] Using TLS")
		log.Fatal(http.ListenAndServeTLS(*listen, *cert, *key, mux))
	} else {
		log.Fatal(http.ListenAndServe(*listen, mux))
	}
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	resp := j{}

	loadTime := time.Now()
	resp["current_time"] = j{
		"string": loadTime,
		"unix":   loadTime.Unix(),
	}

	hostname, _ := os.Hostname()
	resp["hostname"] = hostname

	loadModules(resp, r.URL.Query().Get("modules"))

	respJSON, err := json.Marshal(resp)
	if err != nil {
		log.Fatal("[error] Fatal! Could not construct JSON response: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}
