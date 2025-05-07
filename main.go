package main

import (
	"context"
	"time"

	"Daemon/internal/container"
	"Daemon/internal/docker"
	"Daemon/internal/shared/logger"
)

func main() {
	logger.System("Starting container service test...")

	ctx := context.Background()

	client, err := docker.NewDockerClient()

	if err != nil {
		logger.Error("Failed to initialize Docker client: %v", err)
		return
	}

	service := container.NewService(client)

	// Load egg
	egg, err := container.LoadEgg("nests/minecraft/paper.egg.json")
	if err != nil {
		logger.Error("Error loading egg: %v", err)
		return
	}

	// 1. Create
	ctr, err := service.CreateContainer(ctx, egg)

	if err != nil {
		return
	}

	// 2. Start
	if err := service.StartContainer(ctx, ctr.ID); err != nil {
		return
	}

	// 3. Wait a bit for demo
	time.Sleep(5 * time.Second)

	// 4. Stop
	if err := service.StopContainer(ctx, ctr.ID); err != nil {
		return
	}

	// 5. Get
	retrieved, err := service.GetContainer(ctr.ID)
	if err != nil {
		return
	}
	logger.Info("Retrieved container: %+v", retrieved)

	// 6. Remove
	if err := service.RemoveContainer(ctx, ctr.ID); err != nil {
		return
	}

	logger.System("Container lifecycle test complete.")
}
