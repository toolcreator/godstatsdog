package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	port := flag.Int("port", 8080, "The port godstatsdog listens on")
	updateInterval := flag.Int("interval", 1000, "The update interval in milliseconds")
	flag.Parse()

	statsUpdateTicker := time.NewTicker(time.Duration(*updateInterval) * time.Millisecond)
	go updateMetrics(statsUpdateTicker)

	http.Handle("/metrics", promhttp.Handler())
	err := make(chan error)
	go func() {
		err <- http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	}()

	log.Println(fmt.Sprintf("Listening on port %d.", *port))
	log.Println("Metrics are exposed at /metrics endpoint.")

	<-err
	log.Println(err)
}
