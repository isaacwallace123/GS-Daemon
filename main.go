package main

import (
	"context"

	"Daemon/internal/container"
	"Daemon/internal/docker"
	"Daemon/internal/shared/logger"

	"Daemon/cmd/cli"
	_ "Daemon/cmd/cli/commands"
)

func main() {
	logger.System("🟢 Starting interactive container CLI...")

	ctx := context.Background()

	client, err := docker.NewDockerClient()

	if err != nil {
		logger.Error("Failed to initialize Docker client: %v", err)

		return
	}

	service := container.NewService(client)

	cli.RunShell(ctx, service)
}
