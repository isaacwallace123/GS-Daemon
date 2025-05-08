package metrics

import (
	"Daemon/internal/models"
	"Daemon/internal/shared/logger"
	"context"
	"encoding/json"
	"github.com/docker/docker/client"
)

type Collector struct {
	cli *client.Client
}

func NewCollector(cli *client.Client) *Collector {
	return &Collector{cli: cli}
}

func (collector *Collector) CollectContainerMetrics(ctx context.Context, containerID string) (*models.ContainerMetrics, error) {
	statsResp, err := collector.cli.ContainerStats(ctx, containerID, true)

	if err != nil {
		return nil, logger.Error("Failed to get container stats: %v", err)
	}

	defer statsResp.Body.Close()

	var first, second models.DockerStats
	decoder := json.NewDecoder(statsResp.Body)

	if err := decoder.Decode(&first); err != nil {
		return nil, logger.Error("Failed to decode first stats frame: %v", err)
	}

	if err := decoder.Decode(&second); err != nil {
		return nil, logger.Error("Failed to decode second stats frame: %v", err)
	}

	cpuDelta := float64(second.CPUStats.CPUUsage.TotalUsage - first.CPUStats.CPUUsage.TotalUsage)
	sysDelta := float64(second.CPUStats.SystemUsage - first.CPUStats.SystemUsage)

	cpuPercent := 0.0

	if sysDelta > 0 && cpuDelta > 0 {
		cpuPercent = (cpuDelta / sysDelta) * float64(len(second.CPUStats.CPUUsage.PercpuUsage)) * 100
	}

	memUsedMB := float64(second.MemoryStats.Usage) / 1024 / 1024
	memLimitMB := float64(second.MemoryStats.Limit) / 1024 / 1024
	memPercent := (memUsedMB / memLimitMB) * 100

	rx, tx := 0.0, 0.0

	for _, net := range second.Networks {
		rx += float64(net.RxBytes)
		tx += float64(net.TxBytes)
	}

	return &models.ContainerMetrics{
		ContainerID: containerID,
		Name:        second.Name,
		CPUPercent:  cpuPercent,
		MemUsageMB:  memUsedMB,
		MemPercent:  memPercent,
		NetRxKB:     rx / 1024,
		NetTxKB:     tx / 1024,
	}, nil
}
