package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	port := 8080
	updateInterval := 1000

	statsUpdateTicker := time.NewTicker(time.Duration(updateInterval) * time.Millisecond)
	go updateMetrics(statsUpdateTicker)

	http.Handle("/metrics", promhttp.Handler())
	err := make(chan error)
	go func() {
		err <- http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}()

	log.Println(fmt.Sprintf("Listening on port %d.", port))
	log.Println("Metrics are exposed at /metrics endpoint.")

	<-err
	log.Println(err)
}
