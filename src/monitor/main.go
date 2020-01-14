package main

import (
	"flag"
	"log"
)

var (
	dataDir  = flag.String("data-dir", defaultDataDir(), "Directory with monitoring scripts and data")
	httpAddr = flag.String("http-addr", ":7331", "Address and port for HTTP server")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		switch args[0] {
		case "collect":
			if len(args) < 2 {
				log.Fatalf("Usage: monitor collect <monitoring_item>")
			}
			err := Collect(args[1])
			dieOnError(err)
		case "plot":
			if len(args) < 2 {
				log.Fatalf("Usage: monitor plot <monitoring_item>")
			}
			opts := PlotOptions{SaveAsLatest: true}
			_, err := Plot(args[1], opts)
			dieOnError(err)
		default:
			log.Fatalf("Unknown command: %s", args[0])
		}
		return
	}
	log.Printf("Monitor daemon starting up, HTTP interface on %s", *httpAddr)
	go Cron()
	HttpStartServer(*httpAddr)
}
