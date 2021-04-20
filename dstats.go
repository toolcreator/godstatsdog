package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
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

func parseContainerStats(statsResp *types.ContainerStats) (*types.StatsJSON, error) {
	var ret *types.StatsJSON
	var err error

	statsRespBody, _ := ioutil.ReadAll(statsResp.Body)
	err = statsResp.Body.Close()
	if err == nil {
		err = json.Unmarshal(statsRespBody, &ret)
	}

	return ret, err
}

func calcContainerStats(statsJSON *types.StatsJSON, preTotalCPU, preSystemCPU uint64) dStats {
	// https://docs.docker.com/engine/api/v1.41/#operation/ContainerStats

	var ret dStats

	if statsJSON.Name[0] == '/' {
		statsJSON.Name = statsJSON.Name[1:]
	}
	ret.name = statsJSON.Name
	ret.id = statsJSON.ID

	cpuDelta := float64(statsJSON.CPUStats.CPUUsage.TotalUsage - preTotalCPU)
	sysCPUDelta := float64(statsJSON.CPUStats.SystemUsage - preSystemCPU)
	ret.cpuPercent = float32((cpuDelta / sysCPUDelta) * float64(len(statsJSON.CPUStats.CPUUsage.PercpuUsage)) * 100)

	ret.memUsage = statsJSON.MemoryStats.Usage - statsJSON.MemoryStats.Stats["cache"]
	ret.memLimit = statsJSON.MemoryStats.Limit
	ret.memPercent = float32((float64(ret.memUsage) / float64(ret.memLimit)) * 100)

	ret.netInp = 0
	ret.netOut = 0
	for _, netifStats := range statsJSON.Networks {
		ret.netInp += netifStats.RxBytes
		ret.netOut += netifStats.TxBytes
	}

	ret.blockInp = 0
	ret.blockOut = 0
	for _, blkioStats := range statsJSON.BlkioStats.IoServiceBytesRecursive {
		switch blkioStats.Op {
		case "Read":
			ret.blockInp += blkioStats.Value
		case "Write":
			ret.blockOut += blkioStats.Value
		}
	}

	ret.pids = statsJSON.PidsStats.Current

	return ret
}

var oldStatsJSONs = make(map[string]*types.StatsJSON)

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
			statsJSON, err := parseContainerStats(&statsResp)
			if err != nil {
				log.Println(err)
			} else {
				preTotalCPU := uint64(0)
				preSysCPU := uint64(0)
				oldStats := oldStatsJSONs[container.ID]
				if oldStats != nil {
					preTotalCPU = oldStats.CPUStats.CPUUsage.TotalUsage
					preSysCPU = oldStats.CPUStats.SystemUsage
				}
				ret = append(ret, calcContainerStats(statsJSON, preTotalCPU, preSysCPU))
				oldStatsJSONs[container.ID] = statsJSON
			}
		}
	}

	// cleanup old stats
	var obsoleteIds []string
	for containerID := range oldStatsJSONs {
		found := false
		for _, container := range containers {
			found = container.ID == containerID
			if found {
				break
			}
		}
		if !found {
			obsoleteIds = append(obsoleteIds, containerID)
		}
	}
	for _, containerID := range obsoleteIds {
		delete(oldStatsJSONs, containerID)
	}

	dockerCli.Close()

	return ret, nil
}
