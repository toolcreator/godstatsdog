package main

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var labelNames = [...]string{"container_id", "container_name"}

var (
	cpuPercent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "godstatsdog_cpu_percent",
		Help: "The percentage of the host's CPU the container is using",
	}, labelNames[:])
	memUsage = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "godstatsdog_memory_usage_bytes",
		Help: "The total amount of memory the container is using",
	}, labelNames[:])
	memLimit = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "godstatsdog_memory_limit_bytes",
		Help: "The total amount of memory the container is allowed to use",
	}, labelNames[:])
	memPercent = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "godstatsdog_memory_percent",
		Help: "The percentage of the host's memory the container is using",
	}, labelNames[:])
	netInp = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "godstatsdog_network_received_bytes",
		Help: "The amount of data the container has received over its network interface",
	}, labelNames[:])
	netOut = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "godstatsdog_network_transmitted_bytes",
		Help: "The amount of data the container has transmitted over its network interface",
	}, labelNames[:])
	blockInp = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "godstatsdog_block_read_bytes",
		Help: "The amount of data the container has read from block devices on the host",
	}, labelNames[:])
	blockOut = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "godstatsdog_block_written_bytes",
		Help: "The amount of data the container has written to block devices on the host",
	}, labelNames[:])
	pids = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "godstatsdog_process_ids",
		Help: "The number of processes or threads the container has created",
	}, labelNames[:])
)

func updateMetric(gaugeVec *prometheus.GaugeVec, id string, name string, value float64) {
	gauge, _ := gaugeVec.GetMetricWithLabelValues(id, name)
	gauge.Set(value)
}

func resetMetrics() {
	cpuPercent.Reset()
	memUsage.Reset()
	memLimit.Reset()
	memPercent.Reset()
	netInp.Reset()
	netOut.Reset()
	blockInp.Reset()
	blockOut.Reset()
	pids.Reset()
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
				updateMetric(memUsage, stat.id, stat.name, float64(stat.memUsage))
				updateMetric(memLimit, stat.id, stat.name, float64(stat.memLimit))
				updateMetric(memPercent, stat.id, stat.name, float64(stat.memPercent))
				updateMetric(netInp, stat.id, stat.name, float64(stat.netInp))
				updateMetric(netOut, stat.id, stat.name, float64(stat.netOut))
				updateMetric(blockInp, stat.id, stat.name, float64(stat.blockInp))
				updateMetric(blockOut, stat.id, stat.name, float64(stat.blockOut))
				updateMetric(pids, stat.id, stat.name, float64(stat.pids))
			}
		}
	}
}
