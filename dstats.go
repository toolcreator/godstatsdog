package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type dStats struct {
	id         string
	name       string
	cpuPercent float32
	memUsage   uint64
	memLimit   uint64
	memPercent float32
	netInp     uint64
	netOut     uint64
	blockInp   uint64
	blockOut   uint64
	pids       uint64
}

func parseContainerStats(statsResp types.ContainerStats) (dStats, error) {
	var ret dStats
	var err error

	var stats types.StatsJSON
	decoder := json.NewDecoder(statsResp.Body)
	err = decoder.Decode(&stats)
	if err == nil {
		// https://docs.docker.com/engine/api/v1.41/#operation/ContainerStats

		if stats.Name[0] == '/' {
			stats.Name = stats.Name[1:]
		}
		ret.name = stats.Name
		ret.id = stats.ID

		cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
		sysCPUDelta := float64(stats.CPUStats.SystemUsage - stats.PreCPUStats.SystemUsage)
		ret.cpuPercent = float32((cpuDelta / sysCPUDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100)

		ret.memUsage = stats.MemoryStats.Usage - stats.MemoryStats.Stats["cache"]
		ret.memLimit = stats.MemoryStats.Limit
		ret.memPercent = float32((float64(ret.memUsage) / float64(ret.memLimit)) * 100)

		ret.netInp = 0
		ret.netOut = 0
		for _, netifStats := range stats.Networks {
			ret.netInp += netifStats.RxBytes
			ret.netOut += netifStats.TxBytes
		}

		ret.blockInp = 0
		ret.blockOut = 0
		for _, blkioStats := range stats.BlkioStats.IoServiceBytesRecursive {
			switch blkioStats.Op {
			case "Read":
				ret.blockInp += blkioStats.Value
			case "Write":
				ret.blockOut += blkioStats.Value
			}
		}

		ret.pids = stats.PidsStats.Current
	}
	return ret, err
}

func getDStats() ([]dStats, error) {
	var ret []dStats

	dockerCli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	if err != nil {
		return ret, err
	}

	ctx := context.Background()

	containers, err := dockerCli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return ret, err
	}

	for _, container := range containers {
		statsResp, err := dockerCli.ContainerStatsOneShot(ctx, container.ID)
		if err != nil {
			log.Println(err)
		} else {
			var stats dStats
			stats, err = parseContainerStats(statsResp)
			if err != nil {
				log.Println(err)
			} else {
				ret = append(ret, stats)
			}
		}
	}

	return ret, nil
}
