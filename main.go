package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var labelNames = [...]string{"container_id", "container_name"}

var (
	cpuPercent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "godstatsdog_cpu_percent",
		Help: "The percentage of the host's CPU the container is using",
	}, labelNames[:])
)

func resetMetrics() {
	cpuPercent.Reset()
	// TODO other metrics
}

func updateMetric(gaugeVec *prometheus.GaugeVec, id string, name string, value float64) {
	gauge, _ := gaugeVec.GetMetricWithLabelValues(id, name)
	gauge.Set(value)
}

func updateMetrics(ticker *time.Ticker) {
	for {
		<-ticker.C
		stats, err := getDStats()
		if err != nil {
			log.Println(err)
		} else {
			resetMetrics()
			for _, stat := range stats {
				updateMetric(cpuPercent, stat.id, stat.name, float64(stat.cpuPercent))
				// TODO other metrics
			}
		}
	}
}

func main() {
	port := 8080
	updateInterval := 1000

	statsUpdateTicker := time.NewTicker(time.Duration(updateInterval) * time.Millisecond)
	go updateMetrics(statsUpdateTicker)

	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	log.Println(err)
}
