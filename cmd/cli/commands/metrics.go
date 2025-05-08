package commands

import (
	"Daemon/internal/container"
	"context"

	"Daemon/cmd/cli"
	"Daemon/internal/docker"
	"Daemon/internal/metrics"
	"Daemon/internal/shared/logger"
)

func init() {
	cli.RegisterCommand(&cli.Command{
		Name:        "metrics",
		Description: "Display live container metrics (CPU, RAM, network)",
		Args: []cli.ArgSpec{
			{Name: "containerName", Required: true},
		},
		Execute: runMetrics,
	})
}

func runMetrics(ctx context.Context, _ *container.Service, args []string) error {
	name := args[0]

	client, err := docker.NewDockerClient()
	if err != nil {
		return logger.Error("Docker client error: %v", err)
	}

	id, err := client.ResolveNameToID(ctx, name)
	if err != nil {
		return logger.Error("Failed to resolve container name: %v", err)
	}

	collector := metrics.NewCollector(client.GetClient())

	data, err := collector.CollectContainerMetrics(ctx, id)
	if err != nil {
		return err
	}

	logger.System("📊 Metrics for %s (%s)", data.Name, data.ContainerID)
	logger.Info("🧠 RAM: %.2f MB (%.1f%%)", data.MemUsageMB, data.MemPercent)
	logger.Info("🔥 CPU: %.2f%%", data.CPUPercent)
	logger.Info("📡 Network: RX %.2f KB / TX %.2f KB", data.NetRxKB, data.NetTxKB)

	return nil
}
