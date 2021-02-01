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
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	log.Println(err)
}
